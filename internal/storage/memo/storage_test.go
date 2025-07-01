package memo_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tehrelt/wm-test/internal/models"
	"github.com/tehrelt/wm-test/internal/storage"
	"github.com/tehrelt/wm-test/internal/storage/memo"
)

func TestSaveAndFind(t *testing.T) {
	st := memo.New()

	task := &models.Task{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		Status:    models.PSProcessing,
	}

	err := st.Save(t.Context(), task)
	if err != nil {
		t.Fatalf("saving error: %s", err.Error())
	}

	got, err := st.Find(t.Context(), task.Id)
	if err != nil {
		t.Fatalf("find error: %s", err.Error())
	}

	if got.Id != task.Id {
		t.Fatalf("id mismatch got(%s) != actual(%s)", got.Id.String(), task.Id.String())
	}

	if got.Status != task.Status {
		t.Fatalf("status mismatch got(%s) != actual(%s)", got.Status.String(), task.Status.String())
	}
}

func TestSaveUpdateFind(t *testing.T) {

	st := memo.New()

	task := &models.Task{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		Status:    models.PSProcessing,
	}
	newStatus := models.PSDone

	err := st.Save(t.Context(), task)
	if err != nil {
		t.Fatalf("saving error: %s", err.Error())
	}

	got, err := st.Find(t.Context(), task.Id)
	if err != nil {
		t.Fatalf("find error: %s", err.Error())
	}

	got.Status = newStatus

	if err := st.Save(t.Context(), got); err != nil {
		t.Fatalf("updating error: %s", err.Error())
	}

	got, err = st.Find(t.Context(), task.Id)
	if err != nil {
		t.Fatalf("find error: %s", err.Error())
	}

	if got.Status != newStatus {
		t.Fatalf("status not updated. mismatch got(%s) != actual(%s)", got.Status.String(), newStatus)
	}

}

func TestNotFound(t *testing.T) {
	st := memo.New()

	notExistId := uuid.New()

	_, err := st.Find(t.Context(), notExistId)
	if err == nil {
		t.Fatalf("expected error")
	}

	if !errors.Is(err, storage.ErrTaskNotFound) {
		t.Fatalf("invalid error: %s", err.Error())
	}

}

func TestList(t *testing.T) {
	st := memo.New()

	type ID struct {
		id   uuid.UUID
		seen bool
	}

	ids := []ID{}
	for range 3 {
		id := uuid.New()

		task := &models.Task{
			Id:        id,
			Status:    models.PSWait,
			CreatedAt: time.Now(),
		}

		_ = st.Save(t.Context(), task)
		ids = append(ids, ID{id: id})
	}

	tasks := st.List(t.Context())
	if len(tasks) == 0 {
		t.Fatal("tasks are empty")
	}

	contains := func(target uuid.UUID) *ID {
		for i := range ids {
			if target == ids[i].id {
				return &ids[i]
			}
		}
		return nil
	}

	for _, task := range tasks {
		if v := contains(task.Id); v != nil {
			v.seen = true
		}
	}

	fl := true
	for _, rec := range ids {
		if !rec.seen {
			fl = false
		}
	}

	if !fl {
		t.Fatalf("not all ids seen, %v", ids)
	}
}
