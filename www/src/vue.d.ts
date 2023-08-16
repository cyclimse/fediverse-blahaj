import { DefaultService } from "./api/generated"

declare module 'vue' {
    interface ComponentCustomProperties {
        $api: DefaultService
    }
}