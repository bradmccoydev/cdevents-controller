package kubernetes

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

func GetResults(logger *zap.Logger) {
	config := ctrl.GetConfigOrDie()
	dynamic := dynamic.NewForConfigOrDie(config)
	namespace := "k8sgpt-operator-system"

	resourceId := schema.GroupVersionResource{
		Group:    "core.k8sgpt.ai",
		Version:  "v1alpha1",
		Resource: "results",
	}

	resultList, err := dynamic.Resource(resourceId).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("Failed to list Result object:", zap.Error(err))
	}

	for _, item := range resultList.Items {
		name := item.GetName()
		resultData := item.Object

		details, ok := resultData["spec"].(map[string]interface{})["details"].(string)
		if !ok {
			logger.Info("Failed to extract details field from Result object")
		}

		errorArray, ok := resultData["spec"].(map[string]interface{})["error"].([]interface{})
		if !ok || len(errorArray) == 0 {
			logger.Info("Failed to extract error field from Result object")
		}

		errorText, ok := errorArray[0].(map[string]interface{})["text"].(string)
		if !ok {
			logger.Info("Failed to extract text field from Result object")
		}

		fmt.Println("*****")
		fmt.Println("Name:", name)
		fmt.Println("Error Text:", errorText)
		fmt.Println("Details:", details)
		fmt.Println("*****")
	}
}
