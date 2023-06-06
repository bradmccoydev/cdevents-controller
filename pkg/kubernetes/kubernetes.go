package provider

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Result struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResultSpec   `json:"spec,omitempty"`
	Status ResultStatus `json:"status,omitempty"`
}

type ResultStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

type ResultList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Result `json:"items"`
}

type Backend string

type ResultSpec struct {
	Backend      `json:"backend"`
	Kind         string    `json:"kind"`
	Name         string    `json:"name"`
	Error        []Failure `json:"error"`
	Details      string    `json:"details"`
	ParentObject string    `json:"parentObject"`
}

type Failure struct {
	Text      string      `json:"text,omitempty"`
	Sensitive []Sensitive `json:"sensitive,omitempty"`
}

type Sensitive struct {
	Unmasked string `json:"unmasked,omitempty"`
	Masked   string `json:"masked,omitempty"`
}

func GetResults(name string) string {
	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	dynamic := dynamic.NewForConfigOrDie(config)
	namespace := "k8sgpt-operator-system"

	resourceId := schema.GroupVersionResource{
		Group:    "core.k8sgpt.ai",
		Version:  "v1alpha1",
		Resource: "results.core.k8sgpt.ai",
	}

	resultList, err := dynamic.Resource(resourceId).Namespace(namespace).
		Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		// Handle error
		log.Fatalf("Failed to get ResultList: %v", err)
	}

	// Extract the items from the ResultList
	items, found := resultList.Object["items"]
	if !found {
		// Handle if items field not found
		log.Fatalf("Items not found in ResultList")
	}

	// Type assertion to []interface{} since items is an array
	resultListItems, ok := items.([]interface{})
	if !ok {
		// Handle if items is not of the expected type
		log.Fatalf("Failed to convert items to []interface{}")
	}

	// Iterate over the resultListItems
	for _, item := range resultListItems {
		// Type assertion to map[string]interface{} since each item is a map
		resultItem, ok := item.(map[string]interface{})
		if !ok {
			// Handle if item is not of the expected type
			log.Fatalf("Failed to convert item to map[string]interface{}")
		}

		// Extract the desired fields from the resultItem map
		// Adjust the field names as per your CRD schema
		name := resultItem["metadata"].(map[string]interface{})["name"].(string)
		details := resultItem["spec"].(map[string]interface{})["details"].(string)

		// Print the extracted values
		fmt.Printf("Name: %s, Details: %s\n", name, details)
	}

	//x := resultList.DeepCopyInto().Object
	if err != nil {
		fmt.Errorf("Error getting Provider: %s", err)
	}

	//clientset := kubernetes.NewForConfigOrDie(config)
	// secret, err := clientset.CoreV1().Secrets(namespace).
	// 	Get(ctx, name, metav1.GetOptions{})
	// GroupVersion := schema.GroupVersion{Group: "core.k8sgpt.ai", Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	// SchemeBuilder := &scheme.Builder{GroupVersion: GroupVersion}
	// SchemeBuilder.Register(&Result{}, &ResultList{})
	// resultList := &corev1.ResultList{}
	// listOptions := metav1.ListOptions{}

	// err = clientset.RESTClient().
	// 	Get().
	// 	Namespace(namespace).
	// 	Resource("results").
	// 	VersionedParams(&listOptions, metav1.ParameterCodec).
	// 	Do(context.TODO()).
	// 	Into(resultList)
	// if err != nil {
	// 	log.Fatalf("Failed to get ResultList: %v", err)
	// }

	// // Print the retrieved custom resources
	// for _, result := range resultList.Items {
	// 	fmt.Printf("Name: %s, Details: %s\n", result.Name, result.Spec.Details)
	// }

	return "hi"
}
