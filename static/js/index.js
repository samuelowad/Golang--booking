function prompt(){
    let toast= function(c){
        const {msg='',icon='success',position='top-end'}=c
        const Toast = Swal.mixin({
            toast: true,
            title:msg,
            position,
            showConfirmButton: false,
            icon,
            timer: 3000,
            timerProgressBar: true,

            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire()

    }

    let success= function(c){
        const {msg='',title='',footer=''}=c

        Swal.fire({
            icon: 'success',
            title,
            text: msg,
            footer
        })
    }

    let error= function(c){
        const {msg='',title='',footer=''}=c

        Swal.fire({
            icon: 'error',
            title,
            text: msg,
            footer
        })
    }

    async function custom(c){
        const {msg='',title='',icon='',showConfirmButton=true}=c


        const { value: result } = await Swal.fire({
            showConfirmButton,
            icon,
            title,
            html:msg,
            focusConfirm: false,
            backdrop: false,
            showCancelButton: true,
            willOpen:()=>{
                if(c.willOpen){
                    c.willOpen()
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value
                ]
            },
            didOpen:()=>{
                if(c.didOpen){
                    c.didOpen()
                }
            }
        })

        if (result) {
            if(result.dismiss!== Swal.DismissReason.cancel){
                if(result.value !==""){
                    if(c.callback !== undefined){
                        c.callback(result)
                    }
                }else{
                    c.callback(false)
                }
            }else{
                c.callback(false)}
        }}

    return {
        toast,
        success,
        error,
        custom
    }
}