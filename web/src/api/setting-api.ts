import API from "./api";
import {ServerError} from "./server-error";

const SettingAPI = {
    getSetting(key: string): Promise<any | ServerError> {
        return API.getObject("/settings/" + key, undefined)
    },
    updateSetting(key: string, value: any): Promise<any | ServerError> {
        return API.postObject("/settings/" + key, value)
    }
}

export default SettingAPI;
