package atgen

var yamlTestFuncPerAPIVersion = `
- name: TestFoo
  routerFunc: getRouter
  apiVersions:
    - v1beta1
    - v1beta2
    - v1
  vars:
    key: val
    foo:
      bar: baz
  tests:
    - path: /{apiVersion}/user
      method: post
      req:
        type: form
        params:
          userId: "1"
          name: John
          foo: true
        headers:
          x-api-key: test
      res:
        status: 201
        params:
          foo: bar
      vars:
         foo: bar
      register: baz
`

var yamlTestPerAPIVersion = `
- name: TestFoo
  routerFunc: getRouter
  tests:
    - apiVersions:
        - v1beta1
        - v1beta2
        - v1
      path: /{apiVersion}/user
      method: post
      req:
        type: form
        params:
          userId: "1"
          name: John
        headers:
          x-api-key: test
      res:
        status: 201
        params:
          foo: bar
    - apiVersions:
        - v1
      path: /{apiVersion}/user
      method: post
      req:
        type: form
        params:
          userId: "1"
          name: John
        headers:
          x-api-key: test
      res:
        status: 201
        params:
          foo: bar
`

var yamlTestFuncAndTestPerAPIVersion = `
- name: TestFoo
  routerFunc: getRouter
  apiVersions:
    - v1beta1
    - v1beta2
    - v1
  tests:
    - apiVersions:
        - v1beta1
        - v1beta2
      path: /{apiVersion}/user
      method: post
      req:
        type: form
        params:
          userId: "1"
          name: John
        headers:
          x-api-key: test
      res:
        status: 201
        params:
          foo: bar
    - apiVersions:
        - v1
      path: /{apiVersion}/user
      method: post
      req:
        type: form
        params:
          userId: "1"
          name: John
        headers:
          x-api-key: test
      res:
        status: 201
        params:
          foo: bar
`

var yamlTestFuncWithSubtests = `
- name: TestWithSubtests
  routerFunc: getRouter
  apiVersions:
    - v1
  tests:
    - path: /{apiVersion}/user
      method: post
      req:
        type: form
        params:
          userId: "1"
          name: John
      res:
        status: 201
    - path: /{apiVersion}/user/1
      method: get
      res:
        status: 200
        params:
          isFoo: true
          isBar: false
    - subtests:
        - name: SubFoo
          tests:
            - path: /{apiVersion}/user/1/foo
              method: delete
              res:
                status: 204
              vars:
                foo: bar
            - path: /{apiVersion}/user/1
              method: get
              res:
                status: 200
                params:
                  isFoo: false
            - path: /{apiVersion}/user/2/foo
              method: delete
              res:
                status: 404
                params:
                  foo: bar
`
