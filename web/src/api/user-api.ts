import API from "./api";
import {ServerError} from "./server-error";

const UserAPI = {
    getUser(id: string, queries: {
        "profile": boolean,
        "achievements": boolean,
        "annual-ranks": boolean
    }): Promise<any | ServerError> {
        return API.getObject("/user/" + id, queries)
    },
    updateUser(id: string, data: {
        "profile": any | undefined,
        "achievements": any[] | undefined,
        "annualRanks": any[] | undefined,
    }): Promise<any | ServerError> {
        return API.postObject("/user/" + id, data)
    },
    listUsers(queries: {
        "limit": number,
        "filter-name": string,
        "filter-class": string,
        "filter-email": string,
        "filter-role": number,
        "belowId": number
    }): Promise<any[] | ServerError> {
        return API.getObject("/users/", queries).then(v => {
            return v instanceof ServerError ? v : v["users"]
        })
    },
    getUserStats(): Promise<any | ServerError> {
        return API.getObject("/user-stats/", undefined)
    },
    uploadProfileCover(data: Blob): Promise<any | ServerError> {
        return API.postBlob("/user-profile-cover/", data)
    }
}

export default UserAPI;
