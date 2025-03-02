function Prompt() {
    let toast = function(c) {
        const {
            msg = "",
            icon = "success",
            position = "top-end",
        } = c;
        const Toast = Swal.mixin({
        toast: true,
        title: msg,
        position: position,
        icon: icon,
        showConfirmButton: false,
        timer: 3000,
        timerProgressBar: true,
        didOpen: (toast) => {
            toast.onmouseenter = Swal.stopTimer;
            toast.onmouseleave = Swal.resumeTimer;
        }
        });
        Toast.fire({});
    }

    let success = function(c) {
        const {
            icon = "success",
            msg = "",
            title = "",
            href = "index.html",
            linktxt = "Link!"
        } = c;
        Swal.fire({
            icon: icon,
            title: title,
            text: msg,
            footer: "<a href=" + href + ">" + linktxt + "</a>"
            });
    }

    async function custom(c){
        const {
            icon = "",
            msg = "",
            title = "",
            showConfirmButton = true,

        } = c;

        const { value: result } = await Swal.fire({
            icon : icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined){
                    c.willOpen();
                }
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            }
        })
        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !=="") {
                    if (c.callback !== undefined) {
                        c.callback(result);
                    }
                } else {
                    c.callbacl(false);
                }
            } else {
                c.callback(false);
            }
        }
    }

    return {
        toast: toast,
        success: success,
        custom: custom,
    }
}

function newContact() {
    Swal.fire({
        title: "Add New Contact",
        html:`<form class="forms">
        <div class="row">
          <div class="col-md-12">
          <h4 class="display-3">Add New Contact</h4>
          </div> <br>
          <div class="row">
            <div class="col-md-6">
              <div class="form-group float-right">
                <div class="form-check">
                  <label class="form-check-label">
                  <input type="checkbox" class="form-check-input">
                  Favorite
                  <i class="input-helper"></i>
                  </label>
                </div>
              </div>
            </div>
          </div>
        </div>

        <p class="card-description">Personal Information</p>
        <div class="row">
          <div class="col-md-6">
            <div class="for-group row">
              <label class="col-sm-3 col-form-label">First Name</label>
              {{with .Form.Errors.Get "first_name"}}
              <label class="text-danger">{{.}}</label>
              {{end}}
              <div class="col-sm-9">
                <input type="text" name="first_name" class="form-control" {{with .Form.Errors.Get "first_name"}} is-invalid {{end}}
                id="first_name" required autocomplete="off"> 
              </div>
            </div>
          </div>
          <div class="col-md-6">
            <div class="for-group row">
              <label class="col-sm-3 col-form-label">Last Name</label>
              {{with .Form.Errors.Get "last_name"}}
              <label class="text-danger">{{.}}</label>
              {{end}}
              <div class="col-sm-9">
                <input type="text" name="last_name" class="form-control" {{with .Form.Errors.Get "last_name"}} is-invalid {{end}}
                id="last_name" required autocomplete="off"> 
              </div>
            </div>
          </div>
        </div><br>

        <div class="row">
          <div class="col-md-6">
            <div class="for-group row">
              <label class="col-sm-3 col-form-label">Job Title</label>
              {{with .Form.Errors.Get "job_title"}}
              <label class="text-danger">{{.}}</label>
              {{end}}
              <div class="col-sm-9">
                <input type="text" name="job_title" class="form-control" id="job_title" required autocomplete="off"> 
              </div>
            </div>
          </div>
          <div class="col-md-6">
            <div class="for-group row">
              <label class="col-sm-3 col-form-label">Email</label>
              {{with .Form.Errors.Get "email"}}
              <label class="text-danger">{{.}}</label>
              {{end}}
              <div class="col-sm-9">
                <input type="text" name="email" class="form-control" {{with .Form.Errors.Get "email"}} is-invalid {{end}}
                id="email" required autocomplete="off"> 
              </div>
            </div>
          </div>
        </div><br>

        <div class="row">
          <div class="col-md-6">
            <div class="for-group row">
              <label class="col-sm-3 col-form-label">Company</label>
              <div class="col-sm-9">
                <div class="input-group">
                  <input type="text" name="company" class="form-control" id="company" autocomplete="off"> 
                  <div class="input-group-append">
                    <button class="btn btn-sm btn-outline-primary" type="button">Search</button>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="col-md-6">
            <div class="for-group row">
              <label class="col-sm-3 col-form-label">Objective</label>
              <div class="col-sm-9">
                <select name="objective" class="form-control" id="objective"> 
                  <option>Networking</option>
                  <option>Interview</option>
                  <option>Research</option>
                  <option>Personal</option>
                  <option>Mentor</option>
                  <option>Mentee</option>
                  <option>Other</option>
                </select>
              </div>
            </div>
          </div>
          
        </div><br>

        <p class="card-description">Phone and Links</p>
        <div class="row">
          <div class="col-md-6">
            <div class="for-group row">
              <label class="col-sm-3 col-form-label">Mobile Phone</label>
              <div class="col-sm-9">
                <input type="text" name="mobile_phone" class="form-control" id="mobile_phone" autocomplete="off"> 
              </div>
            </div>
          </div>
          <div class="col-md-6">
            <div class="for-group row">
              <label class="col-sm-3 col-form-label">Work Phone</label>
              <div class="col-sm-9">
                <input type="text" name="work_phone" class="form-control" id="work_phone" autocomplete="off"> 
              </div>
            </div>
          </div>
        </div><br>

        <div class="row">
          <div class="col-md-6">
            <div class="for-group row">
              <label class="col-sm-3 col-form-label">LinkedIn</label>
              <div class="col-sm-9">
                <input type="text" name="linkedin" class="form-control" id="linkedin" autocomplete="off"> 
              </div>
            </div>
          </div>
          <div class="col-md-6">
            <div class="for-group row">
              <label class="col-sm-3 col-form-label">GitHub</label>
              <div class="col-sm-9">
                <input type="text" name="github" class="form-control" id="github" autocomplete="off"> 
              </div>
            </div>
          </div>
        </div><br>

        <div class="row">
          <div class="col-md-12">
            <div class="for-group row">
              <label class="col-sm-2 col-form-label">Description</label>
              <div class="col-sm-12">
                <textarea type="text" name="description" class="form-control" id="description" autocomplete="off" rows="2"></textarea>
              </div>
            </div>
            </div>
          </div><br>

          <div class="row">
          <div class="col-md-12">
            <div class="for-group row">
              <label class="col-sm-2 col-form-label">Notes</label>
              <div class="col-sm-12">
                <textarea type="text" name="notes" class="form-control" id="notes" autocomplete="off" rows="6"></textarea>
              </div>
            </div>
          </div>
        </div><br>
        <button type="submit" class="btn btn-primary mb-2">Submit</button>
      </form>`
    })
}