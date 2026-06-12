Feature: Install ClusterExtension

  As an OLM user I would like to install a cluster extension from catalog
  or get an appropriate information in case of an error.

  Background:
    Given OLM is available
    And an image registry is available

  Scenario:  Install latest available version
    Given a catalog "test" with packages:
      | package | version | channel | replaces | contents                   |
      | test    | 1.2.0   | beta    |          | CRD, Deployment, ConfigMap |
    When ClusterExtension is applied
      """
      apiVersion: olm.operatorframework.io/v1
      kind: ClusterExtension
      metadata:
        name: ${NAME}
      spec:
        namespace: ${TEST_NAMESPACE}
        source:
          sourceType: Catalog
          catalog:
            packageName: ${PACKAGE:test}
            selector:
              matchLabels:
                "olm.operatorframework.io/metadata.name": ${CATALOG:test}
      """
    Then ClusterExtension is rolled out
    And ClusterExtension is available
    And bundle "${PACKAGE:test}.1.2.0" is installed in version "1.2.0"
    And resource "networkpolicy/test-operator-${SCENARIO_ID}-network-policy" is installed
    And resource "configmap/test-configmap-${SCENARIO_ID}" is installed
    And resource "deployment/test-operator-${SCENARIO_ID}" is installed

  @mirrored-registry
  Scenario: Install latest available version from mirrored registry
    Given a catalog "test" with packages:
      | package       | version | channel | replaces | contents                                                                                                                    |
      | test-mirrored | 1.2.0   | beta    |          | CRD, Deployment, ConfigMap, ClusterRegistry(mirrored-registry.operator-controller-e2e.svc.cluster.local:5000) |
    When ClusterExtension is applied
      """
      apiVersion: olm.operatorframework.io/v1
      kind: ClusterExtension
      metadata:
        name: ${NAME}
      spec:
        namespace: ${TEST_NAMESPACE}
        source:
          sourceType: Catalog
          catalog:
            packageName: ${PACKAGE:test-mirrored}
            selector:
              matchLabels:
                "olm.operatorframework.io/metadata.name": ${CATALOG:test}
      """
    Then ClusterExtension is rolled out
    And ClusterExtension is available
    And bundle "${PACKAGE:test-mirrored}.1.2.0" is installed in version "1.2.0"
    And resource "networkpolicy/test-operator-${SCENARIO_ID}-network-policy" is installed
    And resource "configmap/test-configmap-${SCENARIO_ID}" is installed
    And resource "deployment/test-operator-${SCENARIO_ID}" is installed


  Scenario: Report that bundle cannot be installed when it exists in multiple catalogs with same priority
    Given a catalog "test" with packages:
      | package | version | channel | replaces | contents                   |
      | test    | 1.2.0   | beta    |          | CRD, Deployment, ConfigMap |
    And a catalog "extra" with packages:
      | package | version | channel | replaces | contents                   |
      | test    | 1.2.0   | beta    |          | CRD, Deployment, ConfigMap |
    When ClusterExtension is applied
      """
      apiVersion: olm.operatorframework.io/v1
      kind: ClusterExtension
      metadata:
        name: ${NAME}
      spec:
        namespace: ${TEST_NAMESPACE}
        source:
          sourceType: Catalog
          catalog:
            packageName: ${PACKAGE:test}
      """
    Then ClusterExtension reports Progressing as True with Reason Retrying and Message includes:
      """
      found bundles for package "${PACKAGE:test}" in multiple catalogs with the same priority
      """

  @SingleOwnNamespaceInstallSupport
  Scenario: watchNamespace config is required for extension supporting single namespace
    Given a catalog "test" with packages:
      | package                  | version | channel | replaces | contents                                        |
      | single-namespace-operator | 1.0.0   | alpha   |          | CRD, Deployment, InstallMode(SingleNamespace)    |
    And resource is applied
      """
      apiVersion: v1
      kind: Namespace
      metadata:
        name: single-namespace-operator-target
      """
    And ClusterExtension is applied
      """
      apiVersion: olm.operatorframework.io/v1
      kind: ClusterExtension
      metadata:
        name: ${NAME}
      spec:
        namespace: ${TEST_NAMESPACE}
        source:
          sourceType: Catalog
          catalog:
            packageName: ${PACKAGE:test}
            # bundle refers bad image references, so that the deployment never becomes available
            version: 1.0.2
            selector:
              matchLabels:
                "olm.operatorframework.io/metadata.name": ${CATALOG:test}
      """
    Then ClusterObjectSet "${NAME}-1" reports Progressing as False with Reason ProgressDeadlineExceeded
    And ClusterExtension reports Progressing as False with Reason ProgressDeadlineExceeded and Message:
      """
      Revision has not rolled out for 1 minute(s). Last status: Revision 1.0.2 is rolling out.
      """
    And ClusterExtension reports Progressing transition between 1 and 2 minutes since its creation

  @BoxcutterRuntime
  @ProgressDeadline
  Scenario: Report ClusterExtension as not progressing if the rollout does not complete within given timeout
    Given a catalog "test" with packages:
      | package | version | channel | replaces | contents |
      | test    | 1.0.3   | alpha   |          | BadImage |
    And min value for ClusterExtension .spec.progressDeadlineMinutes is set to 1
    And min value for ClusterObjectSet .spec.progressDeadlineMinutes is set to 1
    When ClusterExtension is applied
      """
      apiVersion: olm.operatorframework.io/v1
      kind: ClusterExtension
      metadata:
        name: ${NAME}
      spec:
        namespace: ${TEST_NAMESPACE}
        progressDeadlineMinutes: 1
        source:
          sourceType: Catalog
          catalog:
            packageName: ${PACKAGE:test}
            version: 1.0.3
            selector:
              matchLabels:
                "olm.operatorframework.io/metadata.name": ${CATALOG:test}
      """
    Then ClusterObjectSet "${NAME}-1" reports Progressing as False with Reason ProgressDeadlineExceeded
    And ClusterExtension reports Progressing as False with Reason ProgressDeadlineExceeded and Message:
      """
      Revision has not rolled out for 1 minute(s). Last status: Revision 1.0.3 is rolling out.
      """
    And ClusterExtension reports Progressing transition between 1 and 2 minutes since its creation

  @BoxcutterRuntime
  Scenario:  ClusterObjectSet is annotated with bundle properties
    Given a catalog "test" with packages:
      | package | version | channel | replaces | contents                                                          |
      | test    | 1.2.0   | beta    |          | CRD, Deployment, ConfigMap, Property(olm.test-property=some-value) |
    When ClusterExtension is applied
      """
      apiVersion: olm.operatorframework.io/v1
      kind: ClusterExtension
      metadata:
        name: ${NAME}
      spec:
        namespace: ${TEST_NAMESPACE}
        source:
          sourceType: Catalog
          catalog:
            packageName: ${PACKAGE:test}
            version: 1.2.0
            selector:
              matchLabels:
                "olm.operatorframework.io/metadata.name": ${CATALOG:test}
      """
    # The annotation key and value come from the bundle's metadata/properties.yaml file
    Then ClusterObjectSet "${NAME}-1" contains annotation "olm.properties" with value
      """
      [{"type":"olm.test-property","value":"some-value"}]
      """

  @BoxcutterRuntime
  Scenario: ClusterObjectSet is labeled with owner information
    Given a catalog "test" with packages:
      | package | version | channel | replaces | contents                   |
      | test    | 1.2.0   | beta    |          | CRD, Deployment, ConfigMap |
    When ClusterExtension is applied
      """
      apiVersion: olm.operatorframework.io/v1
      kind: ClusterExtension
      metadata:
        name: ${NAME}
      spec:
        namespace: ${TEST_NAMESPACE}
        source:
          sourceType: Catalog
          catalog:
            packageName: ${PACKAGE:test}
            version: 1.2.0
            selector:
              matchLabels:
                "olm.operatorframework.io/metadata.name": ${CATALOG:test}
      """
    Then ClusterExtension is rolled out
    And ClusterExtension is available
    And ClusterObjectSet "${NAME}-1" has label "olm.operatorframework.io/owner-kind" with value "ClusterExtension"
    And ClusterObjectSet "${NAME}-1" has label "olm.operatorframework.io/owner-name" with value "${NAME}"

  @BoxcutterRuntime
  Scenario: ClusterObjectSet objects are externalized to immutable Secrets
    Given a catalog "test" with packages:
      | package | version | channel | replaces | contents                   |
      | test    | 1.2.0   | beta    |          | CRD, Deployment, ConfigMap |
    When ClusterExtension is applied
      """
      apiVersion: olm.operatorframework.io/v1
      kind: ClusterExtension
      metadata:
        name: ${NAME}
      spec:
        namespace: ${TEST_NAMESPACE}
        source:
          sourceType: Catalog
          catalog:
            packageName: ${PACKAGE:test}
            version: 1.2.0
            selector:
              matchLabels:
                "olm.operatorframework.io/metadata.name": ${CATALOG:test}
      """
    Then ClusterExtension is rolled out
    And ClusterExtension is available
    And ClusterObjectSet "${NAME}-1" phase objects are managed in Kubernetes secrets
    And ClusterObjectSet "${NAME}-1" referred secrets exist in "${OLM_NAMESPACE}" namespace
    And ClusterObjectSet "${NAME}-1" referred secrets are immutable
    And ClusterObjectSet "${NAME}-1" referred secrets contain labels
      | key                                          | value     |
      | olm.operatorframework.io/revision-name       | ${NAME}-1 |
      | olm.operatorframework.io/owner-name          | ${NAME}   |
    And ClusterObjectSet "${NAME}-1" referred secrets are owned by the object set
    And ClusterObjectSet "${NAME}-1" referred secrets have type "olm.operatorframework.io/object-data"

  @DeploymentConfig
  Scenario: deploymentConfig nodeSelector is applied to the operator deployment
    Given a catalog "test" with packages:
      | package | version | channel | replaces | contents                   |
      | test    | 1.2.0   | beta    |          | CRD, Deployment, ConfigMap |
    When ClusterExtension is applied
      """
      apiVersion: olm.operatorframework.io/v1
      kind: ClusterExtension
      metadata:
        name: ${NAME}
      spec:
        namespace: ${TEST_NAMESPACE}
        config:
          configType: Inline
          inline:
            deploymentConfig:
              nodeSelector:
                kubernetes.io/os: linux
        source:
          sourceType: Catalog
          catalog:
            packageName: ${PACKAGE:test}
            selector:
              matchLabels:
                "olm.operatorframework.io/metadata.name": ${CATALOG:test}
      """
    Then resource "deployment/test-operator-${SCENARIO_ID}" matches
      """
      spec:
        template:
          spec:
            nodeSelector:
              kubernetes.io/os: linux
      """

  @BoxcutterRuntime
  Scenario: Install bundle with large CRD
    Given a catalog "test" with packages:
      | package | version | channel | replaces | contents                       |
      | test    | 1.0.0   | beta    |          | LargeCRD(250), Deployment      |
    When ClusterExtension is applied
      """
      apiVersion: olm.operatorframework.io/v1
      kind: ClusterExtension
      metadata:
        name: ${NAME}
      spec:
        namespace: ${TEST_NAMESPACE}
        source:
          sourceType: Catalog
          catalog:
            packageName: ${PACKAGE:test}
            selector:
              matchLabels:
                "olm.operatorframework.io/metadata.name": ${CATALOG:test}
      """
    Then ClusterExtension is rolled out
    And ClusterExtension is available
    And bundle "${PACKAGE:test}.1.0.0" is installed in version "1.0.0"
    And resource "deployment/test-operator-${SCENARIO_ID}" is installed

  Scenario: Applying ClusterExtension with deprecated serviceAccount emits warning
    Given a catalog "test" with packages:
      | package | version | channel | replaces | contents                   |
      | test    | 1.2.0   | beta    |          | CRD, Deployment, ConfigMap |
    When ClusterExtension is applied
      """
      apiVersion: olm.operatorframework.io/v1
      kind: ClusterExtension
      metadata:
        name: ${NAME}
      spec:
        namespace: ${TEST_NAMESPACE}
        serviceAccount:
          name: some-sa
        source:
          sourceType: Catalog
          catalog:
            packageName: ${PACKAGE:test}
            selector:
              matchLabels:
                "olm.operatorframework.io/metadata.name": ${CATALOG:test}
      """
    Then ClusterExtension apply emits warning:
      """
      Warning: spec.serviceAccount is deprecated, ignored, and will be removed in a future release. The operator-controller's cluster-admin service account is used for all cluster interactions.
      """
    And ClusterExtension is rolled out
    And ClusterExtension is available
