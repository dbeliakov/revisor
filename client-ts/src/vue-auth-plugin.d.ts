import {Auth} from "@/auth/auth";

declare module 'vue/types/vue' {
  interface Vue {
    $auth: Auth;
  }
}