<template>
    <div>
        <div v-bind:key="film.Id" v-for="film in films">
            <film-item v-bind:film="film"/>
        </div>
    </div>
</template>
<script>
import FilmItem from "./FilmItem.vue";
export default {
    name: "Films",
    components: {
        FilmItem
    },
    props: ["films"]
}
</script>
<style scoped>

</style>
