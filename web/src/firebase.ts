import { initializeApp } from "firebase/app";
import { getAnalytics } from "firebase/analytics";
import conf from "./conf";

const app = initializeApp(conf.firebase);
const analytics = getAnalytics(app);
