<script>
export default {
	props: ['msg', "timer"],
	emits: ['timerEnd'],
	data: function () {
		return {
			timerData: null,
		};
	},
	watch: {
		timer: {
			handler(value) {
				this.timerData = value;
			},
			immediate: true
		},
		timerData: {
			handler(value) {
				setTimeout(() => {
					this.timerData--;
				}, 1000);
				if (value == 0) {
					this.$emit('timerEnd');
				}
			},
			immediate: true
		},
	},
}
</script>

<template>
	<div class="alert alert-danger" role="alert">
		{{ msg }}
	</div>
</template>

<style>
</style>
