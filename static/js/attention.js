// static/js/attention.js
function Prompt() {
    return {
        toast: function (c) {
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
            });
            Toast.fire({});
        },
        success: function (c) {
            Swal.fire({
                icon: "success",
                title: c.title || "Success!",
                text: c.text || "",
            });
        },
        error: function (c) {
            Swal.fire({
                icon: "error",
                title: c.title || "Error!",
                text: c.text || "",
            });
        },
        custom: function (c) {
            Swal.fire({
                icon: c.icon || "info",
                title: c.title || "",
                html: c.msg || "",
                showCancelButton: true,
                confirmButtonText: c.confirmButtonText || "OK",
            }).then((result) => {
                if (result.isConfirmed && typeof c.callback === "function") {
                    c.callback(true);
                } else if (typeof c.callback === "function") {
                    c.callback(false);
                }
            });
        },
    };
}
