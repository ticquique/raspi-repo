<template>
    <div>
        <form>
            <input type="text" v-model="title" name="title" placeholder="Add film title">
            <input type="text" v-model="description" name="description" placeholder="Add film description">
            <input type="submit" value="Submit" class="btn">
        </form>
    </div>
</template>
<script>

export default {
    name: "AddFilm",
    data() {
        return {
            title: '',
            description: ''
        }
    },
    methods: {}
}
</script>
<style scoped>
    form{
        display: flex;
    }
    input[type="text"] {
        flex: 10;
        padding: 5px;
    }
    input[type="submit"] {
        flex: 2;
    }
</style>

