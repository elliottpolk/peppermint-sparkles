# Peppermint Sparkles Helper

## Go(lang) Sample

The below pipeline example assumes the application has integrated with **_Peppermint Sparkles_** which handles the retrieval and decryption of the secrets stored. It assumes the only requirement is the location of the **_Peppermint Sparkles_** service, the **encryption token** (i.e. **secret of secrets**), and the environment the app is being deployed to. It also assumes the application provides the app_name to the **_Peppermint Sparkles_** service.

```yaml
---

groups:
- name: super-dope
  jobs:
  - unit-tests
  - deploy

resource_types:
- name: artifactory
  type: docker-image
  source:
    repository: pivotalservices/artifactory-resource

resources:
- name: binary-repo
  type: artifactory
  source:
    endpoint: {{ARTIFACTORY_URI}}
    repository: {{ARTIFACTORY_REPO}}
    regex: "super-dope-v(?<version>[0-9].[0-9].[0-9]).tar.bz2"
    username: {{ARTIFACTORY_USER}}
    password: {{ARTIFACTORY_PASS}}
    skip_ssl_verification: true

- name: git-source
  type: git
  source:
    branch: {{GIT_BRANCH}}
    tag_filter: "*.deploy"
    uri: {{GIT_URI}}
    private_key: {{GIT_PRIVATE_KEY}}
    skip_ssl_verification: true

- name: pcf
  type: cf
  source:
    api: {{PCF_API}}
    username: {{PCF_USER}}
    password: {{PCF_PASS}}
    organization: GSD-CAC-DEV
    space: OA-MONTREAL-CAC-DEV
    skip_cert_check: true

jobs:

# DEV environment
- name: unit-tests
  plan:
  - get: super-dope
    resource: git-source
    trigger: true
  - task: unit
    file: super-dope/ci/tasks/unit_test.yml

- name: deploy
  serial: true
  plan:
  - aggregate:
    - { get: source, resource: git-source, passed: [unit-tests], trigger: true }
    - { get: bin, resource: binary-repo }

  ### LOOK HERE FOR THE SAUCE!
  ### "EXPLODED" VERSION FOR REFERENCE ONLY. CAN BE IN YAML / SCRIPT FILE.
  - task: merge
    params:
      TERM: xterm
      sparkles_token: {{SPARKLES_MAGIC}}
      sparkles_env: {{SPARKLES_ENV}}
      sparkles_addr: {{SPARKLES_ADDR}}
    config:
      platform: linux
      image_resource:
        type: docker-image
        source:
          repository: 10.234.24.211:443/peppermint-sparkles-helper
          insecure_registries: ["10.234.24.211:443"]
      inputs:
      - name: source
      - name: bin
      outputs:
      - name: super-dope
      run:
        path: sh
        args:
        - -exec
        - |
          set -o errexit
          set -o xtrace

          TARGET="super-dope"

          # ensure 'super-dope/build/bin' directory exists
          mkdir -p ${TARGET}/build/bin/

          # generate .profile and .vars files          
          set +x \
            && printf 'source .vars && rm .vars' > ${TARGET}/build/bin/.profile \
            && printf "export SPARKLES_TOKEN=\"${sparkles_token}\"" >> ${TARGET}/build/bin/.vars \
            && printf "export SPARKLES_ADDR=\"${sparkles_addr}\"" >> ${TARGET}/build/bin/.vars \
            && printf "export SPARKLES_ENV=\"${sparkles_env}\"" >> ${TARGET}/build/bin/.vars \
            && set -x

          # merge source and binary into expected repo dir 'super-dope'
          mv source/* ${TARGET}/ && \
          mv bin/* ${TARGET}/build/bin/

  - put: pcf
    params:
      manifest: super-dope/pcf/manifest.yml

```

```go
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"git.platform.manulife.io/go-common/db/creds"
	"git.platform.manulife.io/go-common/log"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/crypto"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/crypto/pgp"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/models"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/service"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

var ErrNoToken = errors.New("no token")

var (
	CfgFlag = &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c", "cfg", "confg"},
		Value:   "app.cfg",
		Usage:   "optional path to config file",
	}

	SparklesAddrFlag = &cli.StringFlag{
		Name:    "sparkles-addr",
		Value:   "http://localhost:9001",
		Usage:   "full address for peppermint-sparkles service",
		EnvVars: []string{"SPARKLES_ADDR"},
	}

	SparklesTokenFlag = &cli.StringFlag{
		Name:    "sparkles-token",
		Usage:   "token to use for decrypting peppermint-sparkles content",
		EnvVars: []string{"SPARKLES_TOKEN"},
	}

	SparklesEnvFlag = &cli.StringFlag{
		Name:    "sparkles-env",
		Value:   "local_dev",
		Usage:   "application environment for pulling peppermint-sparkles content",
		EnvVars: []string{"SPARKLES_ENV"},
	}
)

const (
	//	this should not be provided via external methods from the app to attempt the
	//	prevention of typos or possible abuse
	appName string = "super-dope"

	sparkles string = "peppermint-sparkles"

	appParam string = "app_name"
	envParam string = "env"
)

func getRemotely(ctx *cli.Context) (*Conf, error) {
	// not having the address is fine since the service can be configured locally
	addr := ctx.String(SparklesAddrFlag.Names()[0])
	if len(addr) < 1 {
		return cfg, nil
	}

	// must have a decryption token - should now allow unencrypted secrets
	tok := ctx.String(SparklesTokenFlag.Names()[0])
	if len(tok) < 1 {
		return cfg, ErrNoToken
	}

	//	using the default value here should be fine
	env := ctx.String(SparklesEnvFlag.Names()[0])

	from, err := url.Parse(fmt.Sprintf("%s/%s", strings.TrimSuffix(addr, "/"), strings.TrimPrefix(service.PathSecrets, "/")))
	if err != nil {
		return cfg, errors.Wrapf(err, "unable to parse %s URL", sparkles)
	}
	from.RawQuery = url.Values{appParam: {appName}, envParam: {env}}.Encode()

	s, err := get(from.String())
	if err != nil {
		return cfg, errors.Wrapf(err, "unable to retrieve credentials from the %s service", sparkles)
	}

	cred, err := decrypt(s, &pgp.Crypter{Token: []byte(tok)})
	if err != nil {
		return cfg, errors.Wrapf(err, "unable to decrypt credentials from the %s service", sparkles)
	}

	cfg.creds = cred

	log.Debugf("credentials successfully retrieved from %s service", sparkles)
	return cfg, nil
}

func decrypt(what *models.Secret, c crypto.Crypter) (*creds.Credential, error) {
	log.Debugf("attempting to decrypt content from %s service", sparkles)
	res, err := c.Decrypt([]byte(what.Content))
	if err != nil {
		return nil, errors.Wrap(err, "unable to decrypt secrets content")
	}

	log.Debugf("attempting to parse decrypted content from %s service", sparkles)

	cred := &creds.Credential{}
	if err := json.Unmarshal(res, &cred); err != nil {
		return nil, errors.Wrap(err, "uanble to parse content of secret")
	}

	return cred, nil
}

func get(from string) (*models.Secret, error) {
	log.Debugf("attempting to call %s service at %s", sparkles, from)
	res, err := http.Get(from)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to retrieve secrets from %s", sparkles)
	}
	defer res.Body.Close()

	log.Debugf("attempting to read response from %s service", sparkles)
	in, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read in %s response body", sparkles)
	}

	log.Debugf("checking status code from %s service", sparkles)
	if code := res.StatusCode; code != http.StatusOK {
		return nil, errors.Errorf("%s service responded with status code %d and message %s", sparkles, code, string(in))
	}

	log.Debugf("attempting to parse content from %s service", sparkles)

	secret := &models.Secret{}
	if err := json.Unmarshal(in, &secret); err != nil {
		return nil, errors.Wrapf(err, "unable to parse content retrieved from %s", sparkles)
	}

	return secret, nil
}
```

---

## Java Sample 1

The below pipeline example assumes the application has integrated with **_Peppermint Sparkles_** which handles the retrieval and decryption of the secrets stored. It assumes the only requirement is the location of the **_Peppermint Sparkles_** service, the **encryption token** (i.e. **secret of secrets**), and the environment the app is being deployed to. It also assumes the application provides the app_name to the **_Peppermint Sparkles_** service.

```yaml
---
groups:
- name: super-dope
  jobs:
  - unit-tests
  - deploy

resources:
- name: src-develop
  type: git
  source:
    branch: {{GITLAB_BRANCH}}
    uri: {{GITLAB_SSH_URI}}
    private_key: {{GITLAB_PRIVATE_KEY}}
    skip_ssl_verification: true

jobs:
- name: unit-tests
  public: true
  plan:
  - get: source
    resource: src-develop
    trigger: true
  - task: test
    file: source/ci/tasks/unit-test.yml
    params:
        TERM: xterm

- name: deploy
  public: true
  serial: true
  plan:
  - get: source
    resource: src-develop
    trigger: true

  - aggregate:
    - task: services
      file: source/ci/tasks/services.yml
      params:
        TERM: xterm
        api: {{CF_API}}
        username: {{CF_USER}}
        password: {{CF_PASSWORD}}
        organization: {{CF_ORG}}
        space: {{CF_SPACE}}
        environment: {{CF-MANIFEST-ENV}}

    - task: frontend
      params:
        TERM: xterm
      config:
        platform: linux
        image_resource:
          type: docker-image
          source:
            repository: 10.234.24.211:443/gsd-java-nodev6
            insecure_registries: ["10.234.24.211:443"]
        inputs:
        - name: source
        outputs:
        - name: frontend
        run:
          path: sh
          args:
          - -exec
          - |
            set -o errexit
            set -o xtrace

            TARGET_DIR=${PWD}/frontend
            SRC_DIR=${PWD}/source/
            cd ${SRC_DIR}/frontend

            # Download dependencies for the frontend and build
            npm update \
              && npm run ng build -prod \
              && cd ${SRC_DIR} \
              && cp -R . ${TARGET_DIR}/

    - task: backend
      params:
        TERM: xterm
      config:
        platform: linux
        image_resource:
          type: docker-image
          source:
            repository: 10.234.24.211:443/gsd-java-nodev6
            insecure_registries: ["10.234.24.211:443"]
        inputs:
        - name: source
        outputs:
        - name: backend
        run:
          path: sh
          args:
          - -exec
          - |
            set -o errexit
            set -o xtrace

            TARGET_DIR=${PWD}/backend
            OLD_JAR=${PWD}/old_jar
            mkdir -p ${OLD_JAR}

            SRC_DIR=${PWD}/source/
            cd ${SRC_DIR}

            # build
            chmod +x gradlew \
              && ./gradlew clean build
            sleep 3

            cp -R . ${TARGET_DIR}/

  ### LOOK HERE FOR THE SAUCE!
  ### "EXPLODED" VERSION FOR REFERENCE ONLY. CAN BE IN YAML / SCRIPT FILE.
  - task: repack
    params:
      TERM: xterm
      application: {{CF_APPLICATION}}
      sparkles_token: {{SPARKLES_MAGIC}}
      sparkles_id: {{SPARKLES_ID}}
      sparkles_addr: {{SPARKLES_ADDR}}
    config:
      platform: linux
      image_resource:
        type: docker-image
        source:
          repository: 10.234.24.211:443/peppermint-sparkles-helper
          insecure_registries: ["10.234.24.211:443"]
      inputs:
      - name: backend
      - name: frontend
      outputs:
      - name: repacked
      run:
        path: sh
        args:
        - -exec
        - |
          set -o errexit
          set -o xtrace

          BASE=${PWD}
          FRONTEND=${BASE}/frontend
          BACKEND=${BASE}/backend
          TARGET=${BASE}/repacked
          JAR=${TARGET}/build/libs
          OLD_JAR=${TARGET}/old_jar

          cp -rf ${BACKEND}/. ${TARGET}
          cd ${TARGET} && mkdir -p ${OLD_JAR}
          unzip ${JAR}/${application}-exec.jar -d ${OLD_JAR}/ \
            && cp -rf ${FRONTEND}/src/main/resources/. ${OLD_JAR}/BOOT-INF/classes/ \
            && set +x \
            && printf 'source .vars && rm .vars' > ${OLD_JAR}/.profile \
            && printf "export SPARKLES_TOKEN=${sparkles_token}" >> ${OLD_JAR}/.vars \
            && printf "export SPARKLES_ADDR=${sparkles_addr}" >> ${OLD_JAR}/.vars \
            && printf "export SPARKLES_ENV=${sparkles_env}" >> ${OLD_JAR}/.vars \
            && set -x \
            && cd ${OLD_JAR} \
            && ls -A | fastjar cvfm ${JAR}/${application}-exec.jar META-INF/MANIFEST.MF -@ \
            && cd ${TARGET} \
            && rm -rf ${OLD_JAR}

            # NOTE: fastjar does not currently include hidden files so the `ls -A`
            # exposes the hidden files and pipes (|) the results to fastjar

  - task: stage
    file: source/ci/tasks/stage.yml
    params:
      TERM: xterm
      environment: {{CF-MANIFEST-ENV}}
      api: {{CF_API}}
      username: {{CF_USER}}
      password: {{CF_PASSWORD}}
      organization: {{CF_ORG}}
      space: {{CF_SPACE}}
      application: {{CF_APPLICATION}}

```

**TODO:**

* Include Java code snippet of integration

---

## Java Sample 2

The below pipeline example allows for integration of **_Peppermint Sparkles_** into the pipeline rather than integrating directly into the application. This assumes the secrets stored are done in a method friendly with setting environment vars:

```bash
export FOO="something_secret"
export BAR="something_else_secret"
```

The first part of the _pipeline.yml_ file is the basic boilerplate. Scroll down to the comment **_### LOOK HERE FOR THE SAUCE!_** to see the script that will curl the API, decrypt the secret, and inject into the resulting **_.vars_** file. For the current version of PCF, a _.profile_ file can be "sourced" prior to the app starting. The _.profile_ file will "source" the generated **_.vars_** file, which should contain exports retrieved from the **_Peppermint Sparkles_** service. This method ties the environment vars in the **_.vars_** file to the process ID running the application.

```yaml

---
groups:
- name: super-dope
  jobs:
  - unit-tests
  - deploy

resources:
- name: src-develop
  type: git
  source:
    branch: {{GITLAB_BRANCH}}
    uri: {{GITLAB_SSH_URI}}
    private_key: {{GITLAB_PRIVATE_KEY}}
    skip_ssl_verification: true

jobs:
- name: unit-tests
  public: true
  plan:
  - get: source
    resource: src-develop
    trigger: true
  - task: test
    file: source/ci/tasks/unit-test.yml
    params:
        TERM: xterm

- name: deploy
  public: true
  serial: true
  plan:
  - get: source
    resource: src-develop
    trigger: true

  - aggregate:
    - task: services
      file: source/ci/tasks/services.yml
      params:
        TERM: xterm
        api: {{CF_API}}
        username: {{CF_USER}}
        password: {{CF_PASSWORD}}
        organization: {{CF_ORG}}
        space: {{CF_SPACE}}
        environment: {{CF-MANIFEST-ENV}}

    - task: frontend
      params:
        TERM: xterm
      config:
        platform: linux
        image_resource:
          type: docker-image
          source:
            repository: 10.234.24.211:443/gsd-java-nodev6
            insecure_registries: ["10.234.24.211:443"]
        inputs:
        - name: source
        outputs:
        - name: frontend
        run:
          path: sh
          args:
          - -exec
          - |
            set -o errexit
            set -o xtrace

            TARGET_DIR=${PWD}/frontend
            SRC_DIR=${PWD}/source/
            cd ${SRC_DIR}/frontend

            # Download dependencies for the frontend and build
            npm update \
              && npm run ng build -prod \
              && cd ${SRC_DIR} \
              && cp -R . ${TARGET_DIR}/

    - task: backend
      params:
        TERM: xterm
      config:
        platform: linux
        image_resource:
          type: docker-image
          source:
            repository: 10.234.24.211:443/gsd-java-nodev6
            insecure_registries: ["10.234.24.211:443"]
        inputs:
        - name: source
        outputs:
        - name: backend
        run:
          path: sh
          args:
          - -exec
          - |
            set -o errexit
            set -o xtrace

            TARGET_DIR=${PWD}/backend
            OLD_JAR=${PWD}/old_jar
            mkdir -p ${OLD_JAR}

            SRC_DIR=${PWD}/source/
            cd ${SRC_DIR}

            # build
            chmod +x gradlew \
              && ./gradlew clean build
            sleep 3

            cp -R . ${TARGET_DIR}/

  ### LOOK HERE FOR THE SAUCE!
  ### "EXPLODED" VERSION FOR REFERENCE ONLY. CAN BE IN YAML / SCRIPT FILE.
  - task: repack
    params:
      TERM: xterm
      application: {{CF_APPLICATION}}
      sparkles_token: {{SPARKLES_MAGIC}}
      sparkles_id: {{SPARKLES_ID}}
      sparkles_addr: {{SPARKLES_ADDR}}
    config:
      platform: linux
      image_resource:
        type: docker-image
        source:
          repository: 10.234.24.211:443/peppermint-sparkles-helper
          insecure_registries: ["10.234.24.211:443"]
      inputs:
      - name: backend
      - name: frontend
      outputs:
      - name: repacked
      run:
        path: sh
        args:
        - -exec
        - |
          set -o errexit
          set -o xtrace

          BASE=${PWD}
          FRONTEND=${BASE}/frontend
          BACKEND=${BASE}/backend
          TARGET=${BASE}/repacked
          JAR=${TARGET}/build/libs
          OLD_JAR=${TARGET}/old_jar

          cp -rf ${BACKEND}/. ${TARGET}
          cd ${TARGET} && mkdir -p ${OLD_JAR}
          unzip ${JAR}/${application}-exec.jar -d ${OLD_JAR}/ \
            && cp -rf ${FRONTEND}/src/main/resources/. ${OLD_JAR}/BOOT-INF/classes/ \
            && set +x \
            && printf 'source .vars && rm .vars' > ${OLD_JAR}/.profile \
            && curl -s https://${sparkles_addr}/api/v1/secrets/${sparkles_id} \
              | jq -r '.content' \
              | base64 -d \
              | gpg -d --output ${OLD_JAR}/.vars --passphrase ${sparkles_token} \
            && set -x \
            && cd ${OLD_JAR} \
            && ls -A | fastjar cvfm ${JAR}/${application}-exec.jar META-INF/MANIFEST.MF -@ \
            && cd ${TARGET} \
            && rm -rf ${OLD_JAR}

            # NOTE: fastjar does not currently include hidden files so the `ls -A`
            # exposes the hidden files and pipes (|) the results to fastjar

  - task: stage
    file: source/ci/tasks/stage.yml
    params:
      TERM: xterm
      environment: {{CF-MANIFEST-ENV}}
      api: {{CF_API}}
      username: {{CF_USER}}
      password: {{CF_PASSWORD}}
      organization: {{CF_ORG}}
      space: {{CF_SPACE}}
      application: {{CF_APPLICATION}}

```
