<template>
  <v-app>
    <v-app-bar app color="primary" dark>
      <div class="d-flex align-center">
        Chrome History Report
      </div>
    </v-app-bar>

    <v-main>
      <HelloWorld v-if="userIsLoggedIn" :password="password" />
      <TheLogin v-else @loginResult="handleLoginResult" />
    </v-main>
  </v-app>
</template>

<script>
import HelloWorld from './components/HelloWorld';
import TheLogin from './components/TheLogin.vue';

export default {
  name: 'App',

  components: {
    HelloWorld,
    TheLogin
  },

  data: () => ({
    userIsLoggedIn: false,
    password: null
  }),
  methods: {
    handleLoginResult(r) {
      this.userIsLoggedIn = r.loginResult
      this.password = r.password
    }
  }
};
</script>
