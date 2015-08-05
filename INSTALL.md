# Installing zedlist


### Download precompiled binaries   
This is the easiest and recomended to use. All binaries have been signed. We use github release, please visit here [release](release) and select the version of your choice.


## Install from source
This will require a golang SDK  installed in your machine. Then go get the project.

		$ go get github.com/gernest/zedlist

This will create the zedlist binary in the $GOPATH/bin directory.


# Usage
After you have acquired the zedlist binary. Then follow these steps to kickstat the application

## Step 1 Create a database
You need to create a postgresql database, there is no any specific of the name or the user. Just be sure you have a valid postgresql connection URL.

The connection string can be of the two forms. Note that whichever you choose its okay.

* The key value style

		dbname=zedlist_test host=localhost port=5432 user=postgres password=postgres sslmode=disable
		
* The URL style

		postgres://postgres:postgres@localhost:5432/zedlist_test?sslmode=disable




### Step 2 Configure Zedlist
Zedlist uses environment variables for configurations the following are the Environmental variables used by zedlist.

#### The following are environment variables that configure zedlist

Naem 						|	Usage				|	Default
-----------------------------|------------------------|-----------------
CONFIG_NAME					| the name of the app	| Zedlist
CONFIG_PORT					| server's port			| 8090'
CONFIG_APPURL				| address of the server  | http://localhost:8090"
CONFIG_MINIMUMAGE			| minimum age for users  | 18
CONFIG_BIRTHDATEFORMAT		| birth date format      | 2 January, 2006
CONFIG_DEFAULTLANG			| default language		| en
CONFIG_DBDIALECT				| type of database(RDMS) | postgres
CONFIG_DBCONN				| database connection URL|postgres://postgres:postgres@localhost/zedlist?sslmode=disable

The most important is the database connection variable `CONFIG_POSTGRESS_CONN`. You **must** set this before running zedlist.

So if you are in a bash shell. You can run this before runnning zedlist command.

		export CONFIG_POSTGRESS_CONN=postgres://postgres:postgres@localhost:5432/zedlist_test?sslmode=disable`
		
		# Then run zedlist commands

Remember to replace the cconnection string with your own.


After you have configured zedlist then you can run migrations. Zedlist ships with a helper for migrations,It is practical if you are running zedlist for the first time to set the admin user. So, run the following command for migrations.

		$ zedlist migrate

### Step 3
Run the application server

		$ zedlist server

### Step 4
There is no step 4. Oh! wait, please click the star button at the top of this page, zedlist is saving lives so we need the star for christ sake.
