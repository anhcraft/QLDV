import {getAuth, GoogleAuthProvider, signInWithPopup, User} from "firebase/auth";
import Cookies from "js-cookie";

const provider = new GoogleAuthProvider();
provider.addScope('https://www.googleapis.com/auth/userinfo.email');

let currentUser: User | undefined = undefined
let userToken: string | undefined = undefined
let initialized = false

const auth = {
    init: function (callback: any) {
        getAuth().onIdTokenChanged(function(user) {
            if (user !== null) {
                getAuth().currentUser?.getIdToken().then((token: string) => {
                    currentUser = user
                    userToken = token
                    if (!initialized) {
                        callback.call(null)
                        initialized = true
                    }
                }, () => {
                    if (!initialized) {
                        callback.call(null)
                        initialized = true
                    }
                });
            } else {
                currentUser = undefined
                if (!initialized) {
                    callback.call(null)
                    initialized = true
                }
            }
        })
    },
    requestAuth: function (onSuccess: any, onError: any) {
        signInWithPopup(getAuth(), provider)
            .then((result) => {
                onSuccess.call(null, result)
            }, (e) => {
                onError.call(null, e)
            })
    },
    getToken: function (): string | undefined {
        return userToken
    },
    setAuthenticated: function (b: boolean) {
        if(b) {
            Cookies.set('qldvauth', "2022", {expires: 7})
        } else {
            Cookies.remove('qldvauth')
        }
    },
    isLoggedIn: function () {
        return Cookies.get('qldvauth') !== undefined && currentUser !== undefined && userToken !== undefined
    }
};

export default auth;
