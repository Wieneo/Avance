import _Vue from "vue"
import Axios from "axios"
import Noty from "noty"

export function Utils<AxiosPlugOptions>(Vue: typeof _Vue): void {
    Vue.prototype.$GetCookie = (Name: string) => {
        const name = Name + '='
        const decodedCookie = decodeURIComponent(document.cookie)
        const ca = decodedCookie.split(';')
        for (let i = 0; i < ca.length; i++) {
            let c = ca[i]
            while (c.charAt(0) === ' ') {
                c = c.substring(1)
            }
            if (c.indexOf(name) === 0) {
                return c.substring(name.length, c.length)
            }
        }
        return ''
    }
    Vue.prototype.$SetCookie = (Name: string, Value: string, Expiration: number) => {
        const d = new Date()
        d.setTime(d.getTime() + (Expiration * 1000))
        const expires = 'expires=' + d.toUTCString()
        document.cookie = Name + '=' + Value + ';' + expires + ';path=/'
    }
    Vue.prototype.$GetRequest = async (URL: string) => {
        const resp = await new Promise((resolve) => {
            const req = new XMLHttpRequest();
            req.open("GET", URL)
            req.onreadystatechange = function(evt){
                if (req.readyState == 4){
                    try{
                        const obj = JSON.parse(req.response)
                        if (obj.Error != undefined){
                            Vue.prototype.$NotifyError(obj.Error)

                            if (obj.Error === "You are currently not authorized!"){
                                //Session timed out
                                window.location.href = "/login"
                            }
                        }
                        resolve(obj)
                    }catch(Exception){
                        console.log(req)
                        resolve(null)
                    }
                }
            }
            req.send()
        })

        return resp
    }
    Vue.prototype.$PostRequest = async (URL: string, Data: any) => {
        const resp = await new Promise((resolve) => {
            const req = new XMLHttpRequest();
            req.open("POST", URL)
            req.onreadystatechange = function(evt){
                if (req.readyState == 4){
                    try{
                        const obj = JSON.parse(req.response)
                        if (obj.Error != undefined){
                            Vue.prototype.$NotifyError(obj.Error)

                            if (obj.Error === "You are currently not authorized!"){
                                //Session timed out
                                window.location.href = "/login"
                            }
                        }
                        resolve(obj)
                    }catch(Exception){
                        console.log(Exception)
                        resolve(null)
                    }
                }
            }
            req.send(JSON.stringify(Data))
        })

        return resp
    }
    Vue.prototype.$Request = async (Method: string, URL: string, Data: any) => {
        const resp = await new Promise((resolve) => {
            const req = new XMLHttpRequest();
            req.open(Method, URL)
            req.onreadystatechange = function(evt){
                if (req.readyState == 4){
                    try{
                        const obj = JSON.parse(req.response)
                        if (obj.Error != undefined){
                            Vue.prototype.$NotifyError(obj.Error)

                            if (obj.Error === "You are currently not authorized!"){
                                //Session timed out
                                window.location.href = "/login"
                            }
                        }
                        resolve(obj)
                    }catch(Exception){
                        console.log(Exception)
                        resolve(null)
                    }
                }
            }
            req.send(JSON.stringify(Data))
        })

        return resp
    }
    Vue.prototype.$NotifySuccess = (Message: string) => {
        new Noty({
            type: "success",
            theme: 'metroui',
            text: Message,
            timeout: 2500
        }).show();
    }
    Vue.prototype.$NotifyError = (Message: string) => {
        new Noty({
            type: "error",
            theme: 'metroui',
            text: Message,
            timeout: 2500
        }).show();
    }
}