package synceddoc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	utils2 "github.com/emirpasic/gods/utils"
	"sync"

	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type Interval struct {
	Start, End int
}

func intervalComparator(a, b interface{}) int {
	aIn := a.(Interval)
	bIn := b.(Interval)

	if aIn.Start < bIn.Start {
		return -1
	}

	if aIn.Start > bIn.Start {
		return 1
	}

	return 0
}

type LogEntry struct {
	PeerId utils.UUID

	LogState map[string] *treemap.Map
}

// LogEntryGob is LogState copy. the only difference is LogState type
type LogEntryGob struct {
	PeerId string

	LogState LogStateGob
}

// LogStateGob is LogState copy. treemap.Map can't be encoded using gob,
// so store it as an array of intervals.
type LogStateGob map[string] []Interval

type SyncedDocument struct {
	id     utils.UUID
	siteId utils.UUID

	cursorPosition      int
	peerCursorPositions map[utils.UUID] int

	localDocument    crdt.Document
	log              []LogEntry
	peerOps          *treemap.Map
	lastLocalOpIndex int

	onChange         ChangeListener
	onPeerConnect    PeerConnectedListener
	onPeerDisconnect PeerDisconnectedListener

	killed bool
	mu     sync.Mutex
}

func (d *SyncedDocument) ConnectSignals(changeListener ChangeListener,
	peerConnectedListener PeerConnectedListener,
	peerDisconnectedListener PeerDisconnectedListener) {

	d.setListeners(changeListener, peerConnectedListener, peerDisconnectedListener)
}

func initDocState(d *SyncedDocument) {
	d.localDocument = crdt.NewBasicDocument(crdt.NewBasicPositionManager(d.siteId))
	d.cursorPosition = 0
	d.peerCursorPositions = make(map[utils.UUID]int)
	d.killed = false

	d.lastLocalOpIndex = 0
	d.peerOps = treemap.NewWith(utils2.StringComparator)
}

func (d *SyncedDocument) setListeners(changeListener ChangeListener,
	peerConnectedListener PeerConnectedListener,
	peerDisconnectedListener PeerDisconnectedListener) {

	setChangeListener(d, changeListener)
	setPeerConnectedListener(d, peerConnectedListener)
	setPeerDisconnectedListener(d, peerDisconnectedListener)
}

func registerTypes() {
	gob.Register(Op{})
	gob.Register(Interval{})
	gob.Register(LogEntryGob{})
	gob.Register(SyncedDocumentState{})
	gob.Register(SyncedDocumentPatch{})

	gob.Register(ConnectRequest{})
	gob.Register(MessageInsert{})
	gob.Register(MessageCRDTInsert{})
	gob.Register(MessageDelete{})
	gob.Register(MessageCRDTDelete{})
	gob.Register(MessagePeerCursor{})
}

func setPeerDisconnectedListener(d *SyncedDocument, listener PeerDisconnectedListener) {
	d.onPeerDisconnect = listener
}

func setPeerConnectedListener(d *SyncedDocument, listener PeerConnectedListener) {
	d.onPeerConnect = listener
}

func setChangeListener(d *SyncedDocument, listener ChangeListener) {
	d.onChange = listener
}

func (d *SyncedDocument) GetLocalOpsFrom(index int) ([]Op, int) {
	res := []Op{}
	lastIndex := 0

	d.mu.Lock()
	defer d.mu.Unlock()

	ops, ok := d.peerOps.Get(string(d.siteId))

	if !ok {
		return res, lastIndex
	}

	localOps := ops.(*treemap.Map)

	it := localOps.Iterator()

	ok = it.Last()
	if ok {
		lastIndex = it.Key().(int)
		for it.Prev() {
			ind := it.Key().(int)
			op := it.Value().(Op)

			if ind <= index {
				break
			}

			res = append(res, op)
		}
	}
	return res, lastIndex
}


// New creates a new, empty document
func New(siteId string) Document {
	doc := new (SyncedDocument)

	doc.id = utils.GenerateNewUUID()
	doc.siteId = utils.UUID(siteId)
	initDocState(doc)
	registerTypes()

	return doc
}

// Open downloads a document having the specified ID
func Open(siteId string, docId string) (Document, error) {
	doc := new (SyncedDocument)

	doc.id = utils.UUID(docId)
	doc.siteId = utils.UUID(siteId)

	initDocState(doc)
	registerTypes()

	return doc, nil
}

// Load deserializes serializedData and creates a document
func Load(siteId string, serializedData []byte) (Document, error) {
	doc := new (SyncedDocument)

	doc.id = utils.GenerateNewUUID()
	doc.siteId = utils.UUID(siteId)
	initDocState(doc)
	registerTypes()

	r := bytes.NewBuffer(serializedData)
	d := gob.NewDecoder(r)

	err := d.Decode(&doc.id)
	if err != nil {
		return nil, err
	}

	err = d.Decode(&doc.lastLocalOpIndex)
	if err != nil {
		return nil, err
	}

	var docLog []LogEntryGob
	err = d.Decode(&docLog)
	if err != nil {
		return nil, err
	}

	for _, gobEntry := range docLog {
		logEntry := LogEntry{}

		logEntry.LogState = gobToLogState(gobEntry.LogState)
		logEntry.PeerId = utils.UUID(gobEntry.PeerId)

		doc.log = append(doc.log, logEntry)
	}

	var gobHistory map[string] map[int] Op

	err = d.Decode(&gobHistory)
	if err != nil {
		return nil, err
	}

	doc.peerOps = gobToPeerOps(gobHistory)

	var documentContent []byte
	err = d.Decode(&documentContent)
	if err != nil {
		return nil, err
	}

	err = doc.localDocument.Deserialize(documentContent)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (d *SyncedDocument) GetID() utils.UUID {
	return d.id
}

func (d *SyncedDocument) GetSiteID() utils.UUID {
	return d.siteId
}

func (d *SyncedDocument) SetChangeListener(listener ChangeListener) {
	setChangeListener(d, listener)
}

func (d *SyncedDocument) SetPeerConnectedListener(listener PeerConnectedListener) {
	setPeerConnectedListener(d, listener)
}

func (d *SyncedDocument) SetPeerDisconnectedListener(listener PeerDisconnectedListener) {
	setPeerDisconnectedListener(d, listener)
}

func (d *SyncedDocument) Serialize() ([]byte, error) {
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)

	err := e.Encode(d.id)
	if err != nil {
		return nil, err
	}

	err = e.Encode(d.lastLocalOpIndex)
	if err != nil {
		return nil, err
	}

	var log []LogEntryGob

	for _, entry := range d.log {
		gobEntry := LogEntryGob{}

		gobEntry.PeerId = string(entry.PeerId)
		gobEntry.LogState = logStateToGob(entry.LogState)

		log = append(log, gobEntry)
	}

	err = e.Encode(log)
	if err != nil {
		return nil, err
	}

	gobHistory := peerOpsToGob(d.peerOps)

	err = e.Encode(gobHistory)
	if err != nil {
		return nil, err
	}

	documentContent, err := d.localDocument.Serialize()
	if err != nil {
		return nil, err
	}

	err = e.Encode(documentContent)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func logStateToGob(state map[string] *treemap.Map) LogStateGob {
	res := make(LogStateGob)

	for peer, intervals := range state {
		var inter []Interval

		it := intervals.Iterator()
		for it.Next() {
			inter = append(inter, it.Key().(Interval))
		}

		res[peer] = inter
	}

	return res
}

func gobToLogState(gobState LogStateGob) map[string] *treemap.Map {
	res := make(map[string] *treemap.Map)

	for peer, intervals := range gobState {
		inter := treemap.NewWith(intervalComparator)

		for _, interval := range intervals {
			inter.Put(interval, nil)
		}

		res[peer] = inter
	}

	return res
}

func (d *SyncedDocument) LocalInsert(index int, val string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	pos := d.localDocument.InsertAtIndex(val, index)

	d.addToPeerOps(Op{
		PeerId:      d.siteId,
		PeerOpIndex: d.lastLocalOpIndex,
		Cmd:         crdt.OpInsert{
			Pos: pos,
			Val: val,
		},
	})

	d.lastLocalOpIndex++
}

func (d *SyncedDocument) RemoteInsert(peerId utils.UUID, peerOpIndex int, position crdt.Position, val string, aux interface{}) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	index := d.localDocument.InsertAtPosition(position, val)

	if index == -1 {
		return false
	}

	if d.onChange != nil {
		d.onChange(MESSAGE_INSERT, MessageInsert{Index: index, Value: val}, aux)
	}

	return true
}

func (d *SyncedDocument) LocalDelete(index int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	pos := d.localDocument.DeleteAtIndex(index)

	d.addToPeerOps(Op{
		PeerId:      d.siteId,
		PeerOpIndex: d.lastLocalOpIndex,
		Cmd:         crdt.OpDelete {
			Pos: pos,
		},
	})

	d.lastLocalOpIndex++
}

func (d *SyncedDocument) RemoteDelete(peerId utils.UUID, peerOpIndex int, position crdt.Position, aux interface{}) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	index := d.localDocument.DeleteAtPosition(position)

	if index == -1 {
		return false
	}

	if d.onChange != nil {
		d.onChange(MESSAGE_DELETE, MessageDelete{Index: index}, aux)
	}

	return true
}

func (d *SyncedDocument) ApplyRemoteOp(op Op, aux interface{}) {

	switch op.Cmd.(type) {
	case crdt.OpInsert:
		crdtOp := op.Cmd.(crdt.OpInsert)
		applied := d.RemoteInsert(op.PeerId, op.PeerOpIndex, crdtOp.Pos, crdtOp.Val, aux)

		if applied {
			d.updatePeerOps(op)
		}
	case crdt.OpDelete:
		crdtOp := op.Cmd.(crdt.OpDelete)
		applied := d.RemoteDelete(op.PeerId, op.PeerOpIndex, crdtOp.Pos, aux)

		if applied {
			d.updatePeerOps(op)
		}
	default:
		fmt.Println("[SyncedDoc] unknown op")
	}
}

func (d *SyncedDocument) updateDocState(peerId utils.UUID, index int) {
	newEntry := LogEntry {
		PeerId:   peerId,
		LogState: make(map[string] *treemap.Map),
	}

	if len(d.log) > 0 {
		// copy interval treemap
		lastEntry := d.log[len(d.log) - 1]
		for peer, intervals := range lastEntry.LogState {
			intervalMap := treemap.NewWith(intervalComparator)

			it := intervals.Iterator()
			for it.Next() {
				intervalMap.Put(it.Key(), it.Value())
			}

			newEntry.LogState[peer] = intervals
		}
	}

	// merge intervals if necessary
	newInterval := Interval{index, index}
	var before, after interface{}
	if state, ok := newEntry.LogState[string(peerId)]; ok {
		before, _ = state.Floor(newInterval)
		after, _ = state.Ceiling(newInterval)
	} else {
		newEntry.LogState[string(peerId)] = treemap.NewWith(intervalComparator)
	}

	if before != nil {
		beforeInterval := before.(Interval)
		if beforeInterval.End == newInterval.Start - 1 {
			newInterval.Start = beforeInterval.Start
		}

		newEntry.LogState[string(peerId)].Remove(beforeInterval)
	}

	if after != nil{
		afterInterval := after.(Interval)

		if newInterval.End == afterInterval.Start - 1 {
			newInterval.End = afterInterval.End
		}

		newEntry.LogState[string(peerId)].Remove(afterInterval)
	}

	newEntry.LogState[string(peerId)].Put(newInterval, nil)

	d.log = append(d.log, newEntry)
}

func (d *SyncedDocument) addToPeerOps(op Op) {
	val, ok := d.peerOps.Get(string(op.PeerId))
	if ok {
		peerOps := val.(*treemap.Map)
		peerOps.Put(op.PeerOpIndex, op)
	} else {
		peerOps := treemap.NewWith(utils2.IntComparator)
		peerOps.Put(op.PeerOpIndex, op)
		d.peerOps.Put(string(op.PeerId), peerOps)
	}

	d.updateDocState(op.PeerId, op.PeerOpIndex)
}

func (d *SyncedDocument) updatePeerOps(op Op) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.addToPeerOps(op)
}

func peerOpsToGob(h *treemap.Map) map[string] map[int] Op {
	res := make(map[string] map[int] Op)

	it := h.Iterator()
	for it.Next() {
		peer := it.Key().(string)
		ops := it.Value().(*treemap.Map)

		opsMap := make(map[int] Op)

		it2 := ops.Iterator()
		for it2.Next() {
			index := it2.Key().(int)
			op := it2.Value().(Op)
			opsMap[index] = op
		}

		res[peer] = opsMap
	}

	return res
}

func gobToPeerOps(h map[string] map[int] Op) *treemap.Map {
	res := treemap.NewWith(utils2.StringComparator)

	for peer, ops := range h {
		opsMap := treemap.NewWith(utils2.IntComparator)

		for index, op := range ops {
			opsMap.Put(index, op)
		}

		res.Put(peer, opsMap)
	}

	return res
}

func (d *SyncedDocument) SetCursor(index int) {
	d.cursorPosition = index
}

func (d *SyncedDocument) Close() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.killed = true

	// free resources
	d.localDocument = nil
}

func (d *SyncedDocument) ToString() string {
	return d.localDocument.ToString()
}

func (d *SyncedDocument) String() string {
	return "[Document " + string(d.id) + "]" + d.localDocument.ToString()
}

func (d *SyncedDocument) GetDocument() crdt.Document {
	return d.localDocument
}


type SyncedDocumentState struct {
	State LogStateGob
}

type SyncedDocumentPatch struct {
	Patch map[string] []Op
}

func (d *SyncedDocument) GetCurrentState() DocumentState {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.getCurrentState()
}

func (d *SyncedDocument) getCurrentState() SyncedDocumentState {
	s := SyncedDocumentState{}

	lastEntry := LogEntry{}
	if len(d.log) > 0 {
		lastEntry = d.log[len(d.log)-1]
	}

	s.State = logStateToGob(lastEntry.LogState)

	return s
}

func GetIntersecting(inter Interval, intervals []Interval) []Interval {
	res := []Interval{}

	for i := 0; i < len(intervals); i++ {
		if inter.Start <= intervals[i].Start && intervals[i].Start <= inter.End ||
			inter.Start <= intervals[i].End && intervals[i].End <= inter.End {
				res = append(res, intervals[i])
		}
	}

	return res
}

func FindMissingIndices(their, our []Interval) []int {
	res := []int{}

	if len(our) == 0 {
		return res
	}

	prev := our[0].Start - 1
	for i := 0; i < len(our); i++ {
		inter := our[i]

		if prev < inter.Start - 1 {
			prev = inter.Start - 1
		}

		intersecting := GetIntersecting(inter, their)

		if len(intersecting) > 0 {
			for _, intersector := range intersecting {
				if intersector.Start > prev {
					// add missing
					for intersector.Start - 1 > prev {
						prev++
						res = append(res, prev)
					}
				}
				prev = intersector.End
			}
		}

		for prev < inter.End {
			prev++
			res = append(res, prev)
		}
	}

	return res
}

func (d *SyncedDocument) CreatePatch(state DocumentState) Patch {
	d.mu.Lock()
	defer d.mu.Unlock()

	s := state.(SyncedDocumentState)
	curState := d.getCurrentState()
	p := SyncedDocumentPatch{}
	p.Patch = make(map[string] []Op)

	for peer, curIntervals := range curState.State {

		intervals, ok := s.State[peer]

		if !ok {
			p.Patch[peer] = []Op{}
			m, ok := d.peerOps.Get(peer)
			if !ok {
				panic ("invalid peer")
			}

			ops := m.(*treemap.Map)
			for _, curInterval := range curIntervals {
				for i := curInterval.Start; i <= curInterval.End; i++ {
					op, ok := ops.Get(i)
					if !ok {
						panic ("invalid op index for peer")
					}
					p.Patch[peer] = append(p.Patch[peer], op.(Op))
				}
			}
		} else {
			missingIndices := FindMissingIndices(intervals, curIntervals)
			if len(missingIndices) > 0 {
				p.Patch[peer] = []Op{}
				m, ok := d.peerOps.Get(peer)
				if !ok {
					panic ("invalid peer")
				}

				ops := m.(*treemap.Map)
				for _, index := range missingIndices {
					op, ok := ops.Get(index)
					if !ok {
						panic ("invalid op index for peer")
					}
					p.Patch[peer] = append(p.Patch[peer], op.(Op))
				}
			}
		}
	}

	return p
}

func (d *SyncedDocument) ApplyPatch(patch Patch) {
	p := patch.(SyncedDocumentPatch)

	for _, ops := range p.Patch {
		for _, op := range ops {
			d.ApplyRemoteOp(op, nil)
		}
	}

}