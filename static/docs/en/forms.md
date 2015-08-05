Forms
=======

Forms are curcisl in any web application, zedlist is not an exception. While searching for a good form library, I tried a different libraries but I ended up picking [gforms]. I piced this because I was already familiar with it, having used it the [aurora] project.

This is actually a primer on using [gforms] lirary, showing examples on how the zedlist forms are composed. As a developer, It will help you understand the forms module of zedlist and hopeful empower you for future contributions.


What is gforms?
=======
According to the project's README, gforms is a flexible forms validation and rendering library for golang web development.

Note that, there are two distinct use, validation and rendering.

Form Validation
======
It is very handy to handle forms manually, where you will have to go the  from parsing the form values from the request, cherry picking values and validate them by hand.

#### using structs for form validation.
gforms support using structs as models for form validation. To do so, first we need to define a struct with the fields we want to validate. This is the  Registration form model for zedlist.

```go
// Register is the registration form
type Register struct {
    FirstName       string    `gforms:"first_name"`
    LastName        string    `gforms:"last_name"`
    MiddleName      string    `gforms:"middle_name"`
    Email           string    `gforms:"email"`
    Password        string    `gforms:"password"`
    ConfirmPassword string    `gforms:"confirm_password"`
    Gender          int       `gforms:"gender"`
    BirthDay        time.Time `gforms:"birth_date"`
}
```

The gforms tags specify the name attribute of the form. gdorms offers a handy `ModelForm`  interface  so we can implement this interface and  be done with validating our registration model.

Zedlist is designed to support multiple languages, for more details about translations please read [translations].

This is the main zedlist `Form` struct  method that implement the the `gforms.ModelForm` interface.

```go
func (f *Form) RegisterForm() gforms.ModelForm {
    var birtdateAttrs = map[string]string{
        "id":    "birthdate",
        "class": "input-large",
    }
    var inputAttrs = map[string]string{
        "class": "input-large",
    }
    return gforms.DefineModelForm(Register{}, gforms.NewFields(
        gforms.NewTextField(
            "first_name",
            gforms.Validators{
                gforms.Required(f.tr.T(msgRequired)),
            },
            gforms.TextInputWidget(inputAttrs),
        ),
        gforms.NewTextField(
            "last_name",
            gforms.Validators{
                gforms.Required(f.tr.T(msgRequired)),
            },
            gforms.TextInputWidget(inputAttrs),
        ),
        gforms.NewTextField(
            "middle_name",
            gforms.Validators{
                gforms.Required(f.tr.T(msgRequired)),
            },
            gforms.TextInputWidget(inputAttrs),
        ),
        gforms.NewTextField(
            "email",
            gforms.Validators{
                gforms.Required(f.tr.T(msgRequired)),
                gforms.EmailValidator(f.tr.T(msgEmail)),
            },
            gforms.BaseTextWidget("email", inputAttrs),
        ),
        gforms.NewTextField(
            "password",
            gforms.Validators{
                gforms.Required(f.tr.T(msgRequired)),
                gforms.MinLengthValidator(6, f.tr.T(msgMinLength, 6)),
            },
            gforms.PasswordInputWidget(inputAttrs),
        ),
        gforms.NewTextField(
            "confirm_password",
            gforms.Validators{
                gforms.Required(f.tr.T(msgRequired)),
                gforms.MinLengthValidator(6, f.tr.T(msgMinLength, 6)),
                EqualValidator{to: "password", Message: f.tr.T(msgEqual)},
            },
            gforms.PasswordInputWidget(inputAttrs),
        ),
        gforms.NewTextField(
            "gender",
            gforms.Validators{
                gforms.Required(f.tr.T(msgRequired)),
            },
            gforms.SelectWidget(
                inputAttrs,
                func() gforms.SelectOptions {
                    return gforms.StringSelectOptions([][]string{
                        {"Select...", "", "true", "false"},
                        {"Male", "0", "false", "false"},
                        {"Female", "1", "false", "false"},
                        {"Zombie", "2", "false", "true"},
                    })
                },
            ),
        ),

        gforms.NewDateTimeField(
            "birth_date",
            settings.App.BirthDateFormat,
            gforms.Validators{
                BirthDateValidator{Limit: settings.App.MinimumAge, Message: f.tr.T(msgAge, settings.App.MinimumAge)},
                gforms.Required(f.tr.T(msgRequired)),
            },
            gforms.BaseTextWidget("text", birtdateAttrs),
        ),
    ))
}
```
