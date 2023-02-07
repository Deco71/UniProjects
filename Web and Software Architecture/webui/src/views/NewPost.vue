<script>
export default {
    emits: ['login'],
    data: function () {
        return {
            errormsg: null,
            loading: false,
            some_data: null,
            photos: null,
			pop: null,
            file: null,
        };
    },
    methods: {
        async upload() {
            this.loading = true;
            this.$axios.post(`/uploadImage`, this.file, {
                headers: {
                    "Content-Type":"image/png",
                },
            }).then(() => {
                this.errormsg = null;
                this.loading = false;
                this.$router.go(-1)
            }).catch((e) => {
                this.loading = false;
                var errore = e.toString();
                if (errore.includes("415")) {
                    this.errormsg = "Unsupported file type\n Please upload a .png file";
                } else {
                    this.errormsg = err.response.data.message;
                }
            })
        },
    },
}
</script>

<template>
    <div class="insert-box">
        <h1>Insert a new post</h1>
        <input type="file" accept="image/png" @change="(event) => file = event.target.files[0]" class="inputfile" id="embedpollfileinput" />
        <button type="button" class="woomy" @click="upload" :disabled="!file">
            Upload
        </button>
        <div v-if=loading class="flex">
				<svg class="feather icon me-2"><use href="/feather-sprite-v4.29.0.svg#loader"/></svg>
				<div class="medium-text">Loading...</div>
		</div>
        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    </div>
</template>

<style>
</style>
