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
        html:``
    })
}