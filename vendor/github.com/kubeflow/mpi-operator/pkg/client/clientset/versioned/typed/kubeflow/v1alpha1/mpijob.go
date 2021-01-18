// Copyright 2018 The Kubeflow Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/kubeflow/mpi-operator/pkg/apis/kubeflow/v1alpha1"
	scheme "github.com/kubeflow/mpi-operator/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MPIJobsGetter has a method to return a MPIJobInterface.
// A group's client should implement this interface.
type MPIJobsGetter interface {
	MPIJobs(namespace string) MPIJobInterface
}

// MPIJobInterface has methods to work with MPIJob resources.
type MPIJobInterface interface {
	Create(*v1alpha1.MPIJob) (*v1alpha1.MPIJob, error)
	Update(*v1alpha1.MPIJob) (*v1alpha1.MPIJob, error)
	UpdateStatus(*v1alpha1.MPIJob) (*v1alpha1.MPIJob, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.MPIJob, error)
	List(opts v1.ListOptions) (*v1alpha1.MPIJobList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MPIJob, err error)
	MPIJobExpansion
}

// mPIJobs implements MPIJobInterface
type mPIJobs struct {
	client rest.Interface
	ns     string
}

// newMPIJobs returns a MPIJobs
func newMPIJobs(c *KubeflowV1alpha1Client, namespace string) *mPIJobs {
	return &mPIJobs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the mPIJob, and returns the corresponding mPIJob object, and an error if there is any.
func (c *mPIJobs) Get(name string, options v1.GetOptions) (result *v1alpha1.MPIJob, err error) {
	result = &v1alpha1.MPIJob{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mpijobs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MPIJobs that match those selectors.
func (c *mPIJobs) List(opts v1.ListOptions) (result *v1alpha1.MPIJobList, err error) {
	result = &v1alpha1.MPIJobList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mpijobs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested mPIJobs.
func (c *mPIJobs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("mpijobs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a mPIJob and creates it.  Returns the server's representation of the mPIJob, and an error, if there is any.
func (c *mPIJobs) Create(mPIJob *v1alpha1.MPIJob) (result *v1alpha1.MPIJob, err error) {
	result = &v1alpha1.MPIJob{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("mpijobs").
		Body(mPIJob).
		Do().
		Into(result)
	return
}

// Update takes the representation of a mPIJob and updates it. Returns the server's representation of the mPIJob, and an error, if there is any.
func (c *mPIJobs) Update(mPIJob *v1alpha1.MPIJob) (result *v1alpha1.MPIJob, err error) {
	result = &v1alpha1.MPIJob{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mpijobs").
		Name(mPIJob.Name).
		Body(mPIJob).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *mPIJobs) UpdateStatus(mPIJob *v1alpha1.MPIJob) (result *v1alpha1.MPIJob, err error) {
	result = &v1alpha1.MPIJob{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mpijobs").
		Name(mPIJob.Name).
		SubResource("status").
		Body(mPIJob).
		Do().
		Into(result)
	return
}

// Delete takes name of the mPIJob and deletes it. Returns an error if one occurs.
func (c *mPIJobs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mpijobs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *mPIJobs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mpijobs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched mPIJob.
func (c *mPIJobs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MPIJob, err error) {
	result = &v1alpha1.MPIJob{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("mpijobs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
