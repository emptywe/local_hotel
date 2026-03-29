function prompt(){
    let toast = function (c){
        const {
            msg = "",
            icon = "success",
            position = "top-end"
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
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let modal = function (c){
        const {
            msg = "",
            title = "",
            footer = "",
            icon = "success"
        } = c;
        Swal.fire({
            icon: icon,
            title: title,
            text: msg,
            footer: footer
        })
    }

    async function custom(c){
        const {
            msg = "",
            title = "",
            icon = "",
            showConfirmButton = true
        } = c;

        const { value: result_dates } = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            backdrop: true,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            }
        })

        if (result_dates){
            if (result_dates.dismiss !== Swal.DismissReason.cancel){
                if (result_dates.value !== ""){
                    if (c.callback !== undefined) {
                        c.callback(result_dates);
                    }
                } else {
                    c.callback(false);
                }
            } else {
                c.callback(false);
            }
        }
    }

    return {
        toast: toast,
        modal: modal,
        custom: custom,
    }
}

function book_room(id,elem, token) {
    elem.addEventListener("click", function (){
        let html = `
    <form id="check-availability-form" action="" method="GET" novalidate class="needs-validation">
        <div class="row">
            <div class="col">
                <div class="row" id="reservation-dates-modal">

                    <div class="col">
                        <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
                    </div>

                    <div class="col">
                        <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure">
                    </div>

                </div>
            </div>
        </div>
    </form>
    `
        attention.custom({
            msg:html,
            title:"Reservation dates",
            willOpen: () => {
                const elem = document.getElementById("reservation-dates-modal");
                const rp = new DateRangePicker(elem, {
                    format:"yyyy-mm-dd",
                    showOnFocus: true,
                    orientation: 'top',
                    autohide: true,
                    minDate: new Date (),
                })
            },

            didOpen: () => {
                document.getElementById('start').removeAttribute('disabled')
                document.getElementById('end').removeAttribute('disabled')
            },

            callback: function (result_dates){

                let form = document.getElementById("check-availability-form");
                let formData = new FormData(form);
                formData.append("csrf_token",token);
                formData.append("room_id", id);

                fetch('/search-availability-json', {
                    method: "post",
                    body: formData,
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.ok) {
                            attention.custom({
                                icon:"success",
                                msg:'<p>Room is available</p>'
                                    + '<p><a href="/book-room?id='
                                    + data.room_id
                                    + '&s='
                                    + data.start_date
                                    + '&e='
                                    + data.end_date
                                    +'"'
                                    + 'class = "btn btn-primary">'
                                    +'Book Now</a></p>',
                                showConfirmButton: false,
                            })
                        } else {
                            attention.modal({
                                msg:"Room is not available",
                                icon:"error",
                            })
                        }
                    })
            }
        })
    });
}