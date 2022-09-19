import API from "./api";
import {ServerError} from "./server-error";

const EventAPI = {
    getEvent(id: string): Promise<any | ServerError> {
        return API.getObject("/event/" + id, undefined)
    },
    updateEvent(id: string, data: any): Promise<any | ServerError> {
        return API.postObject("/event/" + id, data)
    },
    deleteEvent(id: string): Promise<any | ServerError> {
        return API.deleteObject("/event/" + id, {})
    },
    listEvents(queries: {
        "limit": number,
        "below-id": number,
        "begin-date": number,
        "end-date": number
    }): Promise<any[] | ServerError> {
        return API.getObject("/events/", queries).then(v => {
            return v instanceof ServerError ? v : v["events"]
        })
    }
}

export default EventAPI;
