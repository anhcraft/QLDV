import API from "./api";
import {ServerError} from "./server-error";

const UserAPI = {
    getUser(id: any, queries: {
        "profile": boolean,
        "achievements": boolean,
        "annual-ranks": boolean
    }): Promise<any | ServerError> {
        return API.getObject("/user/" + id + Object.entries(queries).map(([key, val]) => `${key}=${val}`).join('&'), {})
    },
    updateUser(id: any, data: {
        "profile": any,
        "achievements": any,
        "annual-ranks": any,
    }): Promise<any | ServerError> {
        return API.postObject("/user/" + id, {})
    },
    listUsers(options: {
        "limit": number,
        "filterName": string,
        "filterClass": string,
        "filterEmail": string,
        "filterRole": number,
        "belowId": number
    }): Promise<any[] | ServerError> {
        return API.getObject("/users/", options).then(v => {
            return v instanceof ServerError ? v : v["users"]
        })
    },
    getUserStats(): Promise<any | ServerError> {
        return API.getObject("/user-stats/", {})
    },
    getUserStats(): Promise<any | ServerError> {
        return API.postObject("/user-stats/", {})
    },
}

export default UserAPI;
