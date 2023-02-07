<script>

export default {
    emits: ["login"],
    data: function () {
        return {
            errormsg: null,
            loading: false,
            some_data: null,

        };
    },
    methods: {
        async refresh() {
            this.errormsg = null;
        },
    },
}
</script>

<template>
	<SearchBar searchStyle="search-mobile"/>
</template>

<style>
</style>
