// Copyright 2021 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	tensorflowv1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1"
	versioned "github.com/kubeflow/tf-operator/pkg/client/clientset/versioned"
	internalinterfaces "github.com/kubeflow/tf-operator/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/kubeflow/tf-operator/pkg/client/listers/tensorflow/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// TFJobInformer provides access to a shared informer and lister for
// TFJobs.
type TFJobInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.TFJobLister
}

type tFJobInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewTFJobInformer constructs a new informer for TFJob type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewTFJobInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredTFJobInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredTFJobInformer constructs a new informer for TFJob type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredTFJobInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KubeflowV1().TFJobs(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KubeflowV1().TFJobs(namespace).Watch(options)
			},
		},
		&tensorflowv1.TFJob{},
		resyncPeriod,
		indexers,
	)
}

func (f *tFJobInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredTFJobInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *tFJobInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&tensorflowv1.TFJob{}, f.defaultInformer)
}

func (f *tFJobInformer) Lister() v1.TFJobLister {
	return v1.NewTFJobLister(f.Informer().GetIndexer())
}
