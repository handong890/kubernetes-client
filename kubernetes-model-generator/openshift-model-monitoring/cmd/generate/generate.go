/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  apimachineryversion "k8s.io/apimachinery/pkg/version"
  "log"
  "reflect"
  "strings"
  "time"
  prometheus "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"

  "os"

  "github.com/fabric8io/kubernetes-client/kubernetes-model-generator/pkg/schemagen"
)

type Schema struct {
  Info                                     apimachineryversion.Info
  APIGroup                                 metav1.APIGroup
  APIGroupList                             metav1.APIGroupList
  BaseKubernetesList                       metav1.List
  ObjectMeta                               metav1.ObjectMeta
  TypeMeta                                 metav1.TypeMeta
  Status                                   metav1.Status
  Patch                                    metav1.Patch
  Time                                     metav1.Time
  Alertmanager                             prometheus.Alertmanager
  AlertmanagerList                         prometheus.AlertmanagerList
  PodMonitor                               prometheus.PodMonitor
  PodMonitorList                           prometheus.PodMonitorList
  Prometheus                               prometheus.Prometheus
  PrometheusList                           prometheus.PrometheusList
  PrometheusRule                           prometheus.PrometheusRule
  PrometheusRuleList                       prometheus.PrometheusRuleList
  ServiceMonitor                           prometheus.ServiceMonitor
  ServiceMonitorList                       prometheus.ServiceMonitorList
}

func main() {
  packages := []schemagen.PackageDescriptor{
    {"k8s.io/api/core/v1", "", "io.fabric8.kubernetes.api.model", "kubernetes_core_"},
    {"k8s.io/apimachinery/pkg/api/resource", "", "io.fabric8.kubernetes.api.model", "kubernetes_resource_"},
    {"k8s.io/apimachinery/pkg/util/intstr", "", "io.fabric8.kubernetes.api.model", "kubernetes_apimachinery_pkg_util_intstr_"},
    {"k8s.io/apimachinery/pkg/runtime", "", "io.fabric8.openshift.api.model.runtime", "kubernetes_apimachinery_pkg_runtime_"},
    {"k8s.io/apimachinery/pkg/version", "", "io.fabric8.kubernetes.api.model.version", "kubernetes_apimachinery_pkg_version_"},
    {"k8s.io/kubernetes/pkg/util", "", "io.fabric8.kubernetes.api.model", "kubernetes_util_"},
    {"k8s.io/kubernetes/pkg/api/errors", "", "io.fabric8.kubernetes.api.model", "kubernetes_errors_"},
    {"k8s.io/kubernetes/pkg/api/unversioned", "", "io.fabric8.kubernetes.api.model", "api_"},
    {"k8s.io/apimachinery/pkg/apis/meta/v1", "", "io.fabric8.kubernetes.api.model", "kubernetes_apimachinery_"},
    {"github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1", "monitoring.coreos.com", "io.fabric8.openshift.api.model.monitoring.v1", "os_monitoring_v1_"},
  }

  typeMap := map[reflect.Type]reflect.Type{
    reflect.TypeOf(time.Time{}): reflect.TypeOf(""),
    reflect.TypeOf(struct{}{}):  reflect.TypeOf(""),
  }
  schema, err := schemagen.GenerateSchema(reflect.TypeOf(Schema{}), packages, typeMap, map[reflect.Type]string{},"monitoring")
  if err != nil {
    fmt.Fprintf(os.Stderr, "An error occurred: %v", err)
    return
  }

  args := os.Args[1:]
  if len(args) < 1 || args[0] != "validation" {
    schema.Resources = nil
  }

  b, err := json.Marshal(&schema)
  if err != nil {
    log.Fatal(err)
  }
  result := string(b)
  result = strings.Replace(result, "\"additionalProperty\":", "\"additionalProperties\":", -1)

  var out bytes.Buffer
  err = json.Indent(&out, []byte(result), "", "  ")
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(out.String())
}
