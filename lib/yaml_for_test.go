package atgen

var yamlTestFuncPerAPIVersion = `
- name: TestFoo
  apiVersions:
    - v1beta1
    - v1beta2
    - v1
  vars:
    adminAPIKey: test
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
