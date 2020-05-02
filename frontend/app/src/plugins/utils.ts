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
        try {
            const Result = await Axios.get(URL, {
                method: 'GET',
                headers: {
                    Authorization: Vue.prototype.$GetCookie('session')
                }
            })
            return Result.data
        } catch (Exception) {
            // 401 means the session key expired
            if (Exception.response != null) {
                return Exception.response.data
            }
            // ToDo: Display error
            // console.log('Unknown error happened: ' + Exception)
            return null
        }
    }
    Vue.prototype.$PostRequest = async (URL: string, Data: any) => {
        try {
            const Result = await Axios.post(URL, Data, {
                method: 'POST',
                headers: {
                    Authorization: Vue.prototype.$GetCookie('session'),
                    'Content-Type': 'application/json'
                }
            })
            return Result.data
        } catch (Exception) {
            // 401 means the session key expired
            if (Exception.response != null) {
                return Exception.response.data
            }
            // ToDo: Display error
            // console.log('Unknown error happened: ' + Exception)
            return null
        }
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