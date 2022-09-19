import lookupErrorCode from "./api";
import API from "./api";

export class ServerError {
    constructor(code: string) {
        this._code = code;
        this._message = API.lookupErrorCode(code);
    }

    get code(): string {
        return this._code;
    }

    get message(): string {
        return this._message;
    }

    private readonly _code: string;
    private readonly _message: string;
}
