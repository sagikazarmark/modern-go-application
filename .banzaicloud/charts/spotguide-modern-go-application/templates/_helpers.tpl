{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "spotguide-modern-go-application.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "spotguide-modern-go-application.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "spotguide-modern-go-application.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Call nested templates.
Source: https://stackoverflow.com/a/52024583/3027614
*/}}
{{- define "call-nested" }}
{{- $dot := index . 0 }}
{{- $subchart := index . 1 }}
{{- $template := index . 2 }}
{{- include $template (dict "Chart" (dict "Name" $subchart) "Values" (index $dot.Values $subchart) "Release" $dot.Release "Capabilities" $dot.Capabilities) }}
{{- end -}}

{{/*
MySQL secret name based on whether an existing secret is provided.
*/}}
{{- define "spotguide-modern-go-application.mysql.userSecretName" -}}
{{- if .Values.mysql.existingUserSecret -}}
{{- .Values.mysql.existingUserSecret -}}
{{- else -}}
{{- printf "%s-mysql" (include "spotguide-modern-go-application.fullname" . | trunc 57 | trimSuffix "-") -}}
{{- end -}}
{{- end -}}

{{/*
MySQL root password secret name based on whether MySQL is being installed or not.
*/}}
{{- define "spotguide-modern-go-application.mysql.rootSecretName" -}}
{{- if .Values.mysql.enabled -}}
{{- include "call-nested" (list . "mysql" "mysql.secretName") -}}
{{- else -}}
{{- required "MySQL (root password) secret is required" .Values.mysql.existingSecret -}}
{{- end -}}
{{- end -}}

{{/*
MySQL host based on whether MySQL is being installed or not.
*/}}
{{- define "spotguide-modern-go-application.mysql.host" -}}
{{- if .Values.mysql.enabled -}}
{{- printf "%s.%s.svc.cluster.local" (include "call-nested" (list . "mysql" "mysql.fullname")) .Release.Namespace -}}
{{- else -}}
{{- required "MySQL host is required!" .Values.mysql.host -}}
{{- end -}}
{{- end -}}

{{/*
MySQL port based on whether MySQL is being installed or not.
*/}}
{{- define "spotguide-modern-go-application.mysql.port" -}}
{{- if .Values.mysql.enabled -}}
{{- .Values.mysql.service.port -}}
{{- else -}}
{{- .Values.mysql.port | default 3306 -}}
{{- end -}}
{{- end -}}

{{/*
Spotguide specific templates.
*/}}
{{- define "repo-tag" }}
{{- if .Values.banzaicloud.organization.name }}
{{- range .Values.banzaicloud.tags }}
{{- if regexMatch "^repo:" . }}
{{- . }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}

{{- define "repo-user" }}
{{- if .Values.banzaicloud.organization.name }}
{{- range .Values.banzaicloud.tags }}
{{- if regexMatch "^repo:" . }}
{{- $repoFullName := regexReplaceAll "^repo:" . "" }}
{{- first (splitList "/" $repoFullName) }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}

{{- define "repo-name" }}
{{- if .Values.banzaicloud.organization.name }}
{{- range .Values.banzaicloud.tags }}
{{- if regexMatch "^repo:" . }}
{{- $repoFullName := regexReplaceAll "^repo:" . "" }}
{{- last (splitList "/" $repoFullName) }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
