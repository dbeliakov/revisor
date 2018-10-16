<template>
    <div class="field">
        <label for="file" class="ui icon button">
            <i class="file icon"></i>
            {{ filename }}</label>
        <input type="file" id="file" style="display:none" v-on:change="updateFile($event)">
    </div>
</template>

<script lang="ts">
import {Component, Vue} from 'vue-property-decorator';

@Component
export default class FileLoader extends Vue {
    public filename: string = 'Добавьте файл';

    public updateFile(event: any) {
        this.$emit('onStartReading');

        const file: File = event.target.files[0];
        if (file.size > 50 * 1024) { // 50 KB
            this.$emit('onReadingError', 'Максимальный размер файла: 50KB');
            this.filename = 'Добавьте файл';
            return;
        }
        const reader: FileReader = new FileReader();
        reader.onloadend = (e) => {
            this.filename = file.name;
            this.$emit('onFinishReading', this.filename, reader.result);
        };
        reader.readAsDataURL(file);
    }
}
</script>