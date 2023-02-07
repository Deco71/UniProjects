<script>
export default {
    props: ["photoId"],
    data: function () {
        return {
            errormsg: null,
            loading: false,
            photos: null,
            blob: null,
            url: null,
        };
    },
    methods: {
        async refresh() {
            this.loading = true;
            await this.$axios.get(`/image/${this.photoId}`, {responseType: 'blob'}).then(response => {
                this.blob = response.data;
                this.blobUrl = URL.createObjectURL(this.blob);
                document.getElementById(`${this.photoId}`).src = this.blobUrl;
                this.loading = false;
            }).catch(err => {
                this.errormsg = err.response.data.message;
                this.loading = false;
            });
        },
    },
    mounted() {
        this.refresh();
    },
    watch: {
    photoId: {
        handler() {
            this.refresh();
        },
        },
    },
}
</script>

<template>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    <LoadingSpinner v-if="loading"></LoadingSpinner>
	<img v-if="!errormsg" :id="photoId">
</template>

<style>
</style>