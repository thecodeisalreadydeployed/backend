package workloadcontroller

import (
	"fmt"
	"github.com/thecodeisalreadydeployed/model"
)

func CreateWorkload(w *model.Payload) {
	fmt.Printf("Workload %s created.", w.Id)
}
