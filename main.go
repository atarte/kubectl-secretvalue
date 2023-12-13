/* Kubectl plugin sample
 *
 * author : Antoine Tarte, tarte.antoine@gmail.com
 * author : Jerome Tarte, jerome.tarte@gmail.com
 */
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Struct defining the arguments and the options of the commands
type Args struct {
	secretName string // Name of the targeted secret
	keyName    string // Key inside the secret that should be retreived
	namespace  string // Target namespace
	help       bool   // help flag. Force the view of command usage if true
	err        error  // err != nil if an error is occurs during computation of CLI arguments and options
}

// Show the usage of the command
func usage() {
	fmt.Println("kubectl secretvalue SECRET_NAME KEY [options]")
	fmt.Println("Give the value of entry represented by the KEY in the Secret")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("    -n, --namespace string    The namespace containting the secret. Default value is the current namespace")
	fmt.Println("    --help                    Show this usage message")
}

// get the value of the key in secret
//
// a Args the Arguments and options collected on CLI
// return:
// string the value of the key in secret,
// err the err that could be raised during processing
func getValue(a Args) (string, error) {
	// setup the kubernetes config from the default ./kube/config file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homeDir, ".kube", "config"))
	if err != nil {
		return "", err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}
	// set the target namespace
	namespace := a.namespace
	if a.namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	secretClient := clientset.CoreV1().Secrets(namespace)
	// Retrieve the target secret
	secret, err := secretClient.Get(context.Background(), a.secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	// Retrieve the value of the key
	value, ok := secret.Data[a.keyName]
	if !ok {
		err = errors.New("Key " + a.keyName + " is missing in Secret " + a.secretName)
		return "", err
	}
	return string(value), nil
}

// Entry point of the program
func main() {
	// get CLI args
	argument := parseArgs(os.Args[1:])
	if argument.help {
		usage()
		os.Exit(0)
	}
	if argument.err != nil {
		fmt.Println(argument.err)
		usage()
		os.Exit(2)
	}
	// get the value
	value, err := getValue(argument)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// show the result
	fmt.Println(value)
}

// Parse the arguments and options collected on CLI
//
// a []string the CLI elemets
//
// return:
// Args the struct containing the collected elements
func parseArgs(a []string) Args {
	i := 0
	r := Args{}
	for i < len(a) {
		switch a[i] {
		case "-n":
			if r.namespace == "" {
				r.namespace = a[i+1]
				i = i + 1
			} else {
				r.err = errors.New("Duplicate option")
			}
		case "--namespace":
			if r.namespace == "" {
				r.namespace = a[i+1]
				i = i + 1
			} else {
				r.err = errors.New("Duplicate option")
			}
		case "--help":
			r.help = true
		default:
			if strings.HasPrefix(a[i], "-") {
				r.err = errors.New("unknow option : " + a[i])
			} else if r.secretName == "" {
				r.secretName = a[i]
			} else if r.keyName == "" {
				r.keyName = a[i]
			} else {
				r.err = errors.New("wrong number of parameters")
			}
		}
		i = i + 1
	}
	if r.secretName == "" || r.keyName == "" {
		r.err = errors.New("wrong number of parameters")
	}
	return r
}
