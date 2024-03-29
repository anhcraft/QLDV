import API from "./api";
import {ServerError} from "./server-error";

const UserAPI = {
    getUser(pid: string, queries: {
        "profile": boolean,
        "achievements": boolean,
        "annual-ranks": boolean
    }): Promise<any | ServerError> {
        return API.getObject("/user/" + pid, queries)
    },
    updateUser(pid: string, data: {
        "profile": any | undefined,
        "achievements": any[] | undefined,
        "annualRanks": any[] | undefined,
    }): Promise<any | ServerError> {
        return API.postObject("/user/" + pid, data)
    },
    listUsers(queries: {
        "limit": number,
        "filter-name": string,
        "filter-class": string,
        "filter-email": string,
        "filter-role": number,
        "below-id": number
    }): Promise<any[] | ServerError> {
        return API.getObject("/users/", queries).then(v => {
            return v instanceof ServerError ? v : v["users"]
        })
    },
    listFeaturedUsers(): Promise<any[] | ServerError> {
        return API.getObject("/users/featured", undefined).then(v => {
            return v instanceof ServerError ? v : v["users"]
        })
    },
    getUserStats(): Promise<any | ServerError> {
        return API.getObject("/user-stats/", undefined)
    },
    uploadProfileCover(data: Blob): Promise<any | ServerError> {
        return API.postBlob("/user-profile-cover/", data)
    },
    uploadProfileAvatar(data: Blob): Promise<any | ServerError> {
        return API.postBlob("/user-profile-avatar/", data)
    }
}

export default UserAPI;
