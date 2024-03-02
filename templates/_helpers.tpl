{{/*
Expand the name of the chart.
*/}}
{{- define "zex-app.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "zex-app.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "zex-app.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "zex-app.labels" -}}
helm.sh/chart: {{ include "zex-app.chart" .context }}
{{ include "zex-app.selectorLabels" (dict "context" .context "component" .component "name" .name) }}
app.kubernetes.io/managed-by: {{ .context.Release.Service }}
app.kubernetes.io/part-of: zex-app
sidecar.istio.io/inject: "true"
{{- with .context.Values.global.additionalLabels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "zex-app.selectorLabels" -}}
{{- if .name -}}
app.kubernetes.io/name: {{ include "zex-app.name" .context }}-{{ .name }}
{{ end -}}
app.kubernetes.io/instance: {{ .context.Release.Name }}
{{- if .component }}
app.kubernetes.io/component: {{ .component }}
{{- end }}
{{- end }}

{{/*
Create frontend name and version as used by the chart label.
*/}}
{{- define "zex-app.frontend.fullname" -}}
{{- printf "%s-%s" (include "zex-app.fullname" .) .Values.frontend.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create backend name and version as used by the chart label.
*/}}
{{- define "zex-app.backend.fullname" -}}
{{- printf "%s-%s" (include "zex-app.fullname" .) .Values.backend.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create database name and version as used by the chart label.
*/}}
{{- define "zex-app.database.fullname" -}}
{{- printf "%s-%s" (include "zex-app.fullname" .) .Values.database.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create the name of the database service account to use
*/}}
{{- define "zex-app.databaseServiceAccountName" -}}
{{- if .Values.frontend.serviceAccount.create -}}
    {{ default (include "zex-app.database.fullname" .) .Values.database.serviceAccount.serviceAccountName }}
{{- else -}}
    {{ default "default" .Values.database.serviceAccount.serviceAccountName }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the backend service account to use
*/}}
{{- define "zex-app.backendServiceAccountName" -}}
    {{- default (include "zex-app.backend.fullname" .) .Values.backend.serviceAccount.serviceAccountName }}
{{- end -}}

{{/*
Create the name of the frontend service account to use
*/}}
{{- define "zex-app.frontendServiceAccountName" -}}
{{- if .Values.frontend.serviceAccount.create -}}
    {{ default (include "zex-app.frontend.fullname" .) .Values.frontend.serviceAccount.serviceAccountName }}
{{- else -}}
    {{ default "default" .Values.frontend.serviceAccount.serviceAccountName }}
{{- end -}}
{{- end -}}

{{/*
Return the appropriate apiVersion for ingress
*/}}
{{- define "zex-app.apiVersion.ingress" -}}
{{- if .Values.apiVersionOverrides.ingress -}}
{{- print .Values.apiVersionOverrides.ingress -}}
{{- else if semverCompare "<1.14-0" (include "zex-app.kubeVersion" $) -}}
{{- print "extensions/v1beta1" -}}
{{- else if semverCompare "<1.19-0" (include "zex-app.kubeVersion" $) -}}
{{- print "networking.k8s.io/v1beta1" -}}
{{- else -}}
{{- print "networking.k8s.io/v1" -}}
{{- end -}}
{{- end -}}

{{/*
Return the target Kubernetes version
*/}}
{{- define "zex-app.kubeVersion" -}}
  {{- default .Capabilities.KubeVersion.Version .Values.kubeVersionOverride }}
{{- end -}}

{{/*
Backend database name
*/}}
{{- define "zex-app.databaseHost" -}}
    {{- if .Values.bitnamiMysql.enabled -}}
        {{- include "mysql.primary.fullname" .Subcharts.mysql -}}
    {{- else if .Values.remoteDb.enabled -}}
        {{- print .Values.remoteDb.host }}
    {{- else -}}
        {{- include "zex-app.database.fullname" . -}}
    {{- end -}}
{{- end -}}

{{/*
Backend database port
*/}}
{{- define "zex-app.databasePort" -}}
    {{- if .Values.bitnamiMysql.enabled -}}
        {{- default "3306" .Values.mysql.primary.service.ports.mysql -}}
    {{- else if .Values.remoteDb.enabled -}}
        {{- default "3306 " .Values.remoteDb.port -}}
    {{- else -}}
        {{- .Values.database.service.mysql.port | default "3306" -}}
    {{- end -}}
{{- end -}}

{{/*Namespace labels*/}}
{{- define "zex-app.namespaceLabelSelector" -}}
{{- if eq .name .context.Values.frontend.name -}}
kubernetes.io/metadata.name: {{ .context.Release.Namespace }}
{{- else -}}
kubernetes.io/metadata.name: {{ include "zex-app.namespace" (dict "name" .name "context" .context) }}
{{- end -}}
{{- end -}}

{{/*
Sub namespaces
*/}}
{{- define  "zex-app.namespace" -}}
{{- printf "%s-%s" .context.Release.Namespace .name }}
{{- end -}}

{{/*Full backend service namespace*/}}
{{- define "zex-app.service.backend" -}}
{{- printf "%s.%s" (include "zex-app.backend.fullname" .) (include "zex-app.namespace" (dict "name" .Values.backend.name "context" .)) }}
{{- end -}}

{{- define "zex-app.service.database" -}}
{{- printf "%s.%s" (include "zex-app.database.fullname" .) (include "zex-app.namespace" (dict "name" .Values.database.name "context" .)) }}
{{- end -}}

{{- define "zex-app.secretsReaderRoleName" -}}
{{- printf "%s:%s" .Release.Namespace "secrets-reader" }}
{{- end -}}