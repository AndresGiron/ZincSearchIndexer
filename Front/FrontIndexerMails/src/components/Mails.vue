<template>
  <div
    class='flex flex-col min-h-screen items-center justify-center min-h-screen from-purple-200 via-purple-300 to-purple-500 bg-gradient-to-br'>
    <h1 class="text-white text-5xl font-bold mb-4">ZINC MAILS
      <font-awesome-icon icon="envelope" class="mr-2" />
    </h1>



    <div v-if="!selectedMail" class="w-4/5 flex mx-10 rounded bg-white mb-4">
      <input class="w-full border-none bg-transparent px-4 py-1 text-gray-400 outline-none focus:outline-none"
        type="search" name="search" placeholder="Search..." v-model="searchQuery" />
      <button type="submit" class="m-2 rounded bg-purple-300 px-4 py-2 text-white" @click="search">
        <svg class="fill-current h-6 w-6" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
          version="1.1" id="Capa_1" x="0px" y="0px" viewBox="0 0 56.966 56.966"
          style="enable-background:new 0 0 56.966 56.966;" xml:space="preserve" width="512px" height="512px">
          <path
            d="M55.146,51.887L41.588,37.786c3.486-4.144,5.396-9.358,5.396-14.786c0-12.682-10.318-23-23-23s-23,10.318-23,23  s10.318,23,23,23c4.761,0,9.298-1.436,13.177-4.162l13.661,14.208c0.571,0.593,1.339,0.92,2.162,0.92  c0.779,0,1.518-0.297,2.079-0.837C56.255,54.982,56.293,53.08,55.146,51.887z M23.984,6c9.374,0,17,7.626,17,17s-7.626,17-17,17  s-17-7.626-17-17S14.61,6,23.984,6z" />
        </svg>
      </button>
    </div>

    <div class="w-4/5 grid grid-cols-12 gap-4">

      <div v-if="!selectedMail" class="w-5/5 overflow-x-auto relative sm:rounded-lg col-span-12">
        <table class="w-full text-sm text-left text-gray-500 dark:text-gray-400">
          <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
            <tr>
              <th scope="col" class="py-3 px-6">Subject</th>
              <th scope="col" class="py-3 px-6">From</th>
              <th scope="col" class="py-3 px-6">To</th>
            </tr>
          </thead>
          <tbody>
            <tr class="bg-white border-b dark:bg-gray-800 dark:border-gray-700" v-for="(email, index) in emails"
              :key="index" @click="showSelectedMail(email)">
              <td class="py-4 px-6">{{ email._source.subject }}</td>
              <td class="py-4 px-6">{{ email._source.from }}</td>
              <td class="py-4 px-6">{{ email._source.to }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="selectedMail"
        class="w-5/5 overflow-x-auto relative sm:rounded-lg col-span-12 disabled:shadow-none ">

        <table class="w-full text-sm text-left text-gray-500 dark:text-gray-400">
          <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
            <tr>
              <th scope="col" class="py-3 px-6">
                <h3>Subject: {{ selectedMail.subject }}</h3>
              </th>
              <th>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr class="bg-white border-b dark:bg-gray-800 dark:border-gray-700">

            </tr>
            <tr class="bg-white border-b dark:bg-gray-800 dark:border-gray-700">
              <td scope="col" class="py-3 px-6">
                <h2> <b>From:</b> {{ selectedMail.from }} <br> <b>To:</b> {{ selectedMail.to }}</h2>
              </td>
              <td scope="col" class="py-3 px-6"></td>
            </tr>
            <tr class="bg-white border-b dark:bg-gray-800 dark:border-gray-700">
              <td scope="col" class="py-3 px-6" style="max-height: 600px;"> {{ selectedMail.body }}</td>
              <td></td>
            </tr>
          </tbody>
        </table>

        <div class="flex justify-center mt-4">
          <button @click="closeModal"
            class="middle none center rounded-lg bg-purple-500 py-3 px-6 font-sans text-xs font-bold uppercase text-white shadow-md shadow-pink-500/20 transition-all hover:shadow-lg hover:shadow-pink-500/40 focus:opacity-[0.85] focus:shadow-none active:opacity-[0.85] active:shadow-none disabled:pointer-events-none disabled:opacity-50 disabled:shadow-none"
            data-ripple-light="true">
            Cerrar
          </button>
        </div>

      </div>



    </div>

    <div v-if="!selectedMail" class="flex justify-center mt-4" center>
      <button class="bg-purple-500 text-white font-bold py-2 px-4 rounded-l"
      :class="{ 'bg-gray-100 text-gray-600 cursor-not-allowed': currentPage === 0 }"
      :disabled="currentPage === 0" 
      @click="prevPage">Prev</button>
      <span class="bg-gray-200 text-gray-700 py-2 px-4">{{ currentPage + 1 }}</span>
      <button class="bg-purple-500 text-white font-bold py-2 px-4 rounded-r"
      :class="{ 'bg-gray-100 text-gray-600 cursor-not-allowed':  emails.length < pageSize }"  
      @click="nextPage">Next</button>
    </div>

  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { library } from '@fortawesome/fontawesome-svg-core';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { faEnvelope } from '@fortawesome/free-solid-svg-icons';

library.add(faEnvelope);


export default defineComponent({
  components:{
    FontAwesomeIcon
  },
  data() {
    return {
      emails: [],
      searchQuery: '',
      currentPage: 0,
      pageSize: 10,
      flagSearch: false,
      selectedMail: null,
    };
  },
  methods: {
    async fetchData() {
      try {
        const requestBody = {
          from: String(this.currentPage * this.pageSize),
          max_results: String(this.pageSize)
        };

        const response = await fetch('http://localhost:9000/all', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(requestBody)
        });

        if (!response.ok) {
          throw new Error('Error al obtener los datos');
        }

        const data = await response.json();
        console.log(data.hits.hits)
        this.emails = data.hits.hits;
        console.log(this.emails, "correos guardados en la variable emails")
      } catch (error) {
        console.error('Error:', error);
      }
    },
    async search() {
      if (this.searchQuery === '') {
        this.currentPage = 0;
        this.fetchData()
      } else {
        try {
          const requestBody = {
            from: String(this.currentPage * this.pageSize),
            max_results: String(this.pageSize),
            term: this.searchQuery + "*"
          };
          const response = await fetch('http://localhost:9000/search', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestBody)
          });

          if (!response.ok) {
            throw new Error('Error al obtener los datos');
          }

          const data = await response.json();
          console.log(data)
          this.emails = data.hits.hits;
          console.log(this.emails.length, "cantidad de correos")


        } catch (error) {
          console.error('Error:', error);
        }
      }
    },
    nextPage() {
      if (this.searchQuery === '') {
        if (this.pageSize <= this.emails.length){
          this.currentPage++;
          this.fetchData();
        }
      } else {
        if (this.pageSize <= this.emails.length){
          this.currentPage++;
          this.search();
        }
      }
    },
    prevPage() {
      if (this.searchQuery === '') {
        if (this.currentPage > 0) {
          this.currentPage--;
          this.fetchData();
        }
      } else {
        if (this.currentPage > 0) {
          this.currentPage--;
          this.search();
        }
      }
    },
    showSelectedMail(email) {
      this.selectedMail = email._source
    },
    closeModal() {
      this.selectedMail = null
    }
  },
  mounted() {
    this.fetchData();
  }
});
</script>


<style scoped>

  .selected-mail-container {
    max-height: 600px;
    overflow-y: auto;
  }

  .selected-mail-container {
    max-height: 200px; /* Ajusta la altura según tus necesidades */
    overflow-y: auto;
  }

</style>
