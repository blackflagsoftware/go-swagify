# Go-swagify
This project is designed to create a swagger document based on the specification for [OpenApi v3.0](https://swagger.io/specification/).  This is done by comments in your Go code.

`Note: currently this does not support every option within the specification, if you want a feature ask for it or create a PR.  See below what is currently supported.`

## Object
Adding comments to your code to produce a swagger docs instead of creating the yaml or json file manually?  Who wouldn't want that?!  The idea would be that on creation or maintenance of you code, make sure the comments match what is desired.  Run this code against the project's directory, Done!  A swagger document is created for you.

## Usage
After compiling this library, add this binary to your path or move into your path.

```
Args:

inputPath: directory to run the parse against
outputPath: directory with file name to save the swagger doc
outputFormat: yaml | json; output format; if omitted, default of 'yaml'
appOutputFormat: json | yaml; your apps' output format; if omitted, default of 'json'
appFieldFormat:  snakeCase | kebabCase | camelCase | pascalCase | lowerCase | upperCase: schema's name field format; if ommitted, default of 'camelCase'
```

If this application is ran without any args the current directory is scanned and the output file is called `swagger.yaml` all other defaults are used.

## Specs
If you are familiar with the OpenApi spec you can specify the object's definition in `components/schemas` and used as a reference for other parts of the spec.  This application relies heavily on that pattern.  In some cases you can specify the object at the level you want, those will be pointed out for you to use.

#### comments
Using the block comment like this:
```
/* go-swagify
*/
```
is a must. Within the block, will depend on the spec you want to create, read on, exmples follow.

#### Schema
Since a lot of the spec is based on a struct with your code.  The parsing of the struct is quite different then the rest, let's start with that.

The struct is vital to golang development so struct tags are leveraged to declare the schema spec.  For each field that you want to include in the schema spec.

```
sw: list of names wanting to associated this field to, delimited by ';'
sw_desc: the description of the field used in the spec
sw_ex: the example to use in the spec

usage:

/* go-swagify
@@struct: User
*/

type User struct {
	Id        int         `db:"id" json:"Id" yaml:"id" sw:"User" sw_desc:"Unique Id for user" sw_ex:"101"`
	FirstName null.String `db:"first_name" json:"FirstName" yaml:"first_name" sw:"User;UserResponse*;UserRequest*" sw_desc:"First name of the user" sw_ex:"John Doe"`
	Age       null.Int    `db:"age" json:"Age" yaml:"age" sw:"User;UserResponse*;UserRequest*" sw_desc:"Age of the user" sw_ex:"42"`
	Active    null.Bool   `db:"active" json:"Active" yaml:"active" sw:"User;UserResponse;UserRequest*" sw_desc:"Is the user's account active" sw_ex:"true | false"`
	CreatedAt null.Time   `db:"created_at" json:"CreatedAt" yaml:"created_at" sw:"User;UserRequest" sw_desc:"Date create for this record" sw_ex:"2020-01-01T00:00:00Z"`
}
```
The `go-swagify` block just tells the parser that you want to inspect the struct called `User`
The `sw:` tag can take list of names (delimited by `'`) if you end the name with `*` then it will make that field for that schema name a required in the spec.

If the `sw` struct is omitted, the field is skipped.

The name is calculated by
- taking the struct tag defined for the field for the `appOutputFormat` (default of `json`)
- if that struct tag is not defined then the field name is formatted by the `altFieldFormat` directive (default `lowerCase`)

With the following struct with the struct tags defined as so, the following `components/schemas` objects would be created:
```
components:
	schemas:
		User:
			type: object
			properties:
				Active:
					type: boolean
					description: Is the user's account active
					example: true | false
				Age:
					type: integer
					description: Age of the user
					example: 42
				CreatedAt:
					type: string
					description: Date create for this record
					example: "2020-01-01T00:00:00Z"
				FirstName:
					type: string
					description: First name of the user
					example: John Doe
				Id:
					type: integer
					description: Unique Id for user
					example: 101
		UserRequest:
			type: object
			required:
			- FirstName
			- Age
			- Active
			properties:
				Active:
					type: boolean
					description: Is the user's account active
					example: true | false
				Age:
					type: integer
					description: Age of the user
					example: 42
				CreatedAt:
					type: string
					description: Date create for this record
					example: "2020-01-01T00:00:00Z"
				FirstName:
					type: string
					description: First name of the user
					example: John Doe
		UserResponse:
			type: object
			required:
			- FirstName
			- Age
			properties:
				Active:
					type: boolean
					description: Is the user's account active
					example: true | false
				Age:
					type: integer
					description: Age of the user
					example: 42
				FirstName:
					type: string
					description: First name of the user
					example: John Doed
```

`NOTE: all examples will show the full yaml path to show where the object is related within the yaml structure`

You can create the schema manually:
```
/* go-swagify
@@schema: UserManual
@@type: object
@@prop_name: Id
@@prop_type: number
@@prop_desc: unique identifier for user
@@prop_ex: 1
@@prop_name: FirstName
@@prop_req: true
@@prop_type: string
@@prop_desc: user's first name
@@prop_ex: John
@@prop_name: Age
@@prop_type: number
@@prop_desc: user's age
@@prop_ex: 42
*/

repeating @@prop_name, @@prop_type, @@prop_desc, @@prop_ex for each field

output:

components:
	schemas:
		UserManual:
			type: object
			required:
			- FirstName
			properties:
				Age:
					type: number
					description: user's age
					example: 42
				FirstName:
					type: string
					description: user's first name
					example: John
				Id:
					type: number
					description: unique identifier for user
					example: 1

note: @@prop_type, @@prop_desc, @@prop_ex and @@prop_req will be related to the proceeding @@prop_name, so order is important
```

You can produce an object of another schema:
```
/* go-swagify
@@schema: UserRef
@@type: object
@@prop_name: data
@@prop_ref: UserManual
*/

components:
	schemas:
		UserRef:
			type: object
			properties:
				data:
				$ref: '#/components/schemas/UserManual'
```

#### Parameter
This will create a spec for the `components/parameters` spec
```
/* go-swagify
@@parameter: UserIndentifierParam
@@name: id
@@in: path (required valid values: query | header | path | cookie)
@@description: Unique User Identifier (optional)
@@required: true (optional: true | false (default))
// schema (optional) see schema.go/SchemaProperty
@@schema_type: number
@@schema_description: User Indentifier
@@schema_example: (optional; not used in this example)
*/

components:
	parameters:
		UserIdentifierParam:
			name: id
			in: path
			description: Unique User identifier
			required: true
			schema:
				type: number
				description: User identifier
```

#### Response & RequestBody
The format for both of these are very similiar, these will fill in `components/responses` and `components/requestBodies` respectfully

```
/* go-swagify
@@response: UserResponseRef
@@desc: the User record
@@content_name: application/json
@@content_ref: UserResponse
*/

/* go-swagify
@@response: UserRequestRef
@@desc: the User record to be created
@@content_name: application/json
@@content_ref: UserRequest
*/

components:
	responses:
		UserResponseRef:
			description: the User record
			content:
				application/json:
				schema:
					$ref: '#/components/schemas/UserResponse'
  	requestBodies:
		UserRequestRef:
			description: the User record to be created
			content:
				application/json:
				schema:
					$ref: '#/components/schemas/UserRequest'
```

#### Path && Operation
Looking at the spec, the path can multiple operations so defining the path name is the key to using both.

Once you define a `path`:
```
/* go-swagify
@@path: /user/{id}
@@parameters.ref: UserIdentifierParam
*/
```
Then you can define as may `operation`s as that go under that path name, like so:
```
/* go-swagify
@@operation: /user/{id}
@@method: get
@@summary: Get User
@@description: Get a single User record by identifier
@@resp_name: 200
@@resp_ref: UserResponseRef
@@resp_name: 400
@@resp_ref: Error
*/

Output would be:

paths:
	/user/{id}:
		parameters:
		- $ref: '#/components/parameters/UserIdentifierParam'
		get:
			summary: Get User
			description: Get a single User record by identifier
			responses:
			"200":
				$ref: '#/components/responses/UserResponseRef'
			"400":
				$ref: '#/components/responses/Error'

note: @@resp_name, @@resp_ref can be repeated as many times as needed but should be the the last lines within the comment block
```

