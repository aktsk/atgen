package atgen

var yamlTestFuncPerAPIVersion = `
- name: TestFoo
  apiVersions:
    - v1beta1
    - v1beta2
    - v1
  vars:
    adminAPIKey: test
    foo:
      bar: baz
  tests:
    - path: /{apiVersion}/money
      method: post
      req:
        params:
          moneyId: "1"
          priority: free
          currency: JPY
          foo: true
        headers:
          x-admin-api-key: test
      res:
        status: 201
        headers:
          location: ""
        params:
          foo: bar
`

var yamlTestPerAPIVersion = `
- name: TestFoo
  tests:
    - apiVersions:
        - v1beta1
        - v1beta2
        - v1
      path: /{apiVersion}/money
      method: post
      req:
        params:
          moneyId: "1"
          priority: free
          currency: JPY
        headers:
          x-admin-api-key: test
      res:
        status: 201
        headers:
          location: ""
        params:
          foo: bar
    - apiVersions:
        - v1
      path: /{apiVersion}/money
      method: post
      req:
        params:
          moneyId: "1"
          priority: free
          currency: JPY
        headers:
          x-admin-api-key: test
      res:
        status: 201
        headers:
          location: ""
        params:
          foo: bar
`

var yamlTestFuncAndTestPerAPIVersion = `
- name: TestFoo
  apiVersions:
    - v1beta1
    - v1beta2
    - v1
  tests:
    - apiVersions:
        - v1beta1
        - v1beta2
      path: /{apiVersion}/money
      method: post
      req:
        params:
          moneyId: "1"
          priority: free
          currency: JPY
        headers:
          x-admin-api-key: test
      res:
        status: 201
        headers:
          location: ""
        params:
          foo: bar
    - apiVersions:
        - v1
      path: /{apiVersion}/user
      method: post
      req:
        params:
          moneyId: "1"
          priority: free
          currency: JPY
        headers:
          x-admin-api-key: test
      res:
        status: 201
        headers:
          location: ""
        params:
          foo: bar
`

var yamlTestFuncWithSubtests = `
- name: TestWithSubtests
  apiVersions:
    - v1
  tests:
    - path: /{apiVersion}/money
      method: post
      req:
        params:
          moneyId: "1"
          priority: free
          currency: JPY
      res:
        status: 201
    - path: /{apiVersion}/money/1
      method: get
      res:
        status: 200
        params:
          isFoo: true
          isBar: false
    - subtests:
        - name: SubFoo
          tests:
            - path: /{apiVersion}/money/1/sub
              method: delete
              res:
                status: 204
            - path: /{apiVersion}/money/1
              method: get
              res:
                status: 200
                params:
                  isFoo: false
            - path: /{apiVersion}/money/2/sub
              method: delete
              res:
                status: 404
                params:
                  foo: bar
`
