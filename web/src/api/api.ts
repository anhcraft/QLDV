import ky, {HTTPError} from "ky";
import conf from "../conf";
import auth from "../auth/auth";
import {ServerError} from "./server-error";

const API = {
    postObject(endpoint: string, data: any): Promise<any | ServerError> {
        const f = (v: string) => {
            return ky.post(conf.server + endpoint, {
                headers: {
                    'content-type': 'application/json',
                    'access-token': v
                },
                body: JSON.stringify(data)
            }).json().then((res: any) => {
                if(res["success"]) {
                    return res["result"]
                } else {
                    return new ServerError(res["error"])
                }
            })
        }
        const token = auth.getToken()
        return token === undefined ? f("") : token.then(f)
    },
    deleteObject(endpoint: string, data: any): Promise<any | ServerError> {
        const f = (v: string) => {
            return ky.delete(conf.server + endpoint, {
                headers: {
                    'content-type': 'application/json',
                    'access-token': v
                },
                body: JSON.stringify(data)
            }).json().then((res: any) => {
                if(res["success"]) {
                    return res["result"]
                } else {
                    return new ServerError(res["error"])
                }
            })
        }
        const token = auth.getToken()
        return token === undefined ? f("") : token.then(f)
    },
    postBlob(endpoint: string, data: Blob): Promise<any | ServerError> {
        const f = (v: string) => {
            return ky.post(conf.server + endpoint, {
                headers: {
                    'content-type': undefined,
                    'access-token': v
                },
                body: data
            }).json().then((res: any) => {
                if(res["success"]) {
                    return res["result"]
                } else {
                    return new ServerError(res["error"])
                }
            })
        }
        const token = auth.getToken()
        return token === undefined ? f("") : token.then(f)
    },
    getObject(endpoint: string, queries: any): Promise<any | ServerError> {
        if(queries !== undefined && Object.keys(queries).length > 0){
            endpoint += "?" + Object.entries(queries).map(([key, val]) => `${key}=${val}`).join('&')
        }
        const f = (v: string) => {
            return ky.get(conf.server + endpoint, {
                headers: {
                    'content-type': 'application/json',
                    'access-token': v
                }
            }).json().then((res: any) => {
                if(res["success"]) {
                    return res["result"]
                } else {
                    return new ServerError(res["error"])
                }
            }, (e) => {
                if(e instanceof HTTPError) {
                    return e.response.json().then(res => {
                        if(res["success"]) {
                            return res["result"]
                        } else {
                            return new ServerError(res["error"])
                        }
                    })
                }
            })
        }
        const token = auth.getToken()
        return token === undefined ? f("") : token.then(f)
    }
}

export default API;
