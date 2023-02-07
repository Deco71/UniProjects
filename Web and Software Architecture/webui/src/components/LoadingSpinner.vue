<script>
export default {
	props: ["load"]
}
</script>

<template>
	<div v-if="load">
		<div style="text-align: center">
			<div class="spinner-border" role="status">
				<span class="visually-hidden">Loading...</span>
			</div>
		</div>
	</div>
	<div v-if="!load"><slot /></div>
</template>

<style></style>
