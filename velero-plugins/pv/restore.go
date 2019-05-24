package pv

import (
	"encoding/json"

	"github.com/fusor/ocp-velero-plugin/velero-plugins/common"
	"github.com/heptio/velero/pkg/plugin/velero"
	"github.com/sirupsen/logrus"
	corev1API "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// RestorePlugin is a restore item action plugin for Velero
type RestorePlugin struct {
	Log logrus.FieldLogger
}

// AppliesTo returns a velero.ResourceSelector that applies to PVs
func (p *RestorePlugin) AppliesTo() (velero.ResourceSelector, error) {
	return velero.ResourceSelector{
		IncludedResources: []string{"persistentvolumes"},
	}, nil
}

// Execute action for the restore plugin for the pv resource
func (p *RestorePlugin) Execute(input *velero.RestoreItemActionExecuteInput) (*velero.RestoreItemActionExecuteOutput, error) {
	p.Log.Info("[pv-restore] Hello from PV RestorePlugin!")

	pv := corev1API.PersistentVolume{}
	itemMarshal, _ := json.Marshal(input.Item)
	json.Unmarshal(itemMarshal, &pv)
	p.Log.Infof("[pv-restore] pv: %s", pv.Name)

	if input.Restore.Annotations[common.MigrateTypeAnnotation] == "copy" {
		p.Log.Infof("[pv-restore] Not a swing PV migration. Skipping pv restore, %s.", pv.Name)
		return velero.NewRestoreItemActionExecuteOutput(input.Item).WithoutRestore(), nil
	}

	var out map[string]interface{}
	objrec, _ := json.Marshal(pv)
	json.Unmarshal(objrec, &out)

	return velero.NewRestoreItemActionExecuteOutput(&unstructured.Unstructured{Object: out}), nil
}