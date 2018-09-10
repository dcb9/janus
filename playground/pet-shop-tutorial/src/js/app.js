App = {
  web3Provider: null,
  contracts: {},
  account: "",

  init: function() {
    // Load pets.
    $.getJSON('../pets.json', function(data) {
      var petsRow = $('#petsRow');
      var petTemplate = $('#petTemplate');

      for (i = 0; i < data.length; i ++) {
        petTemplate.find('.panel-title').text(data[i].name);
        petTemplate.find('img').attr('src', data[i].picture);
        petTemplate.find('.pet-breed').text(data[i].breed);
        petTemplate.find('.pet-age').text(data[i].age);
        petTemplate.find('.pet-location').text(data[i].location);
        petTemplate.find('.btn-adopt').attr('data-id', data[i].id);

        petsRow.append(petTemplate.html());
      }
    });

    App.login()
    return App.initWeb3();
  },

  initWeb3: function() {
    /*
    if (typeof web3 !== 'undefined') {
      App.web3Provider = web3.currentProvider;
    } else {
      // If no injected web3 instance is detected, fall back to Ganache
      App.web3Provider = new Web3.providers.HttpProvider('http://localhost:7545');
    }
    */
    App.web3Provider = new Web3.providers.HttpProvider('http://localhost:23889');
    return App.initContract();
  },

  initContract: function() {
    $.getJSON('Adoption.json', function(data) {
      // Get the necessary contract artifact file and instantiate it with truffle-contract
      var AdoptionArtifact = data;

      App.contracts.Adoption = TruffleContract(AdoptionArtifact);

      // Set the provider for our contract
      App.contracts.Adoption.setProvider(App.web3Provider);

      // Use our contract to retrieve and mark the adopted pets
      return App.markAdopted();
    });

    return App.bindEvents();
  },

  bindEvents: function() {
    $(document).on('click', '.btn-adopt', App.handleAdopt);
  },

  markAdopted: function(adopters, account) {
    var adoptionInstance;
    App.contracts.Adoption.deployed().then(function(instance) {
      adoptionInstance = instance;

      return adoptionInstance.getAdopters.call();
    }).then(function(adopters) {
      for (var i = 0; i < adopters.length; i++) {
        const adopter = adopters[i];
        if (adopter !== '0x0000000000000000000000000000000000000000') {
          $('.panel-pet').eq(i).find('button').text('Adopted').attr('disabled', true);
          $('.panel-pet').eq(i).find('.pet-adopter-container').css('display', 'block');
          let adopterLabel = adopter;
          if (adopter === App.account) {
            adopterLabel = "You"
          }
          $('.panel-pet').eq(i).find('.pet-adopter-address').text(adopterLabel);
        } else {
          $('.panel-pet').eq(i).find('.pet-adopter-container').css('display', 'none');
        }
      }
    }).catch(function(err) {
      console.log(err.message);
    });
  },

  handleAdopt: function(event) {
    event.preventDefault();

    var petId = parseInt($(event.target).data('id'));

    var adoptionInstance;

    App.contracts.Adoption.deployed().then(function(instance) {
      adoptionInstance = instance;

      return adoptionInstance.adopt(petId, {from: App.account});
    }).then(function(result) {
      return App.markAdopted();
    }).catch(function(err) {
      console.log(err.message);
    });
  },

  login: function() {
    let walletAddress = localStorage.getItem("userWalletAddress");
    while (!walletAddress) {
      walletAddress = window.prompt("Please enter your wallet address");
      if (walletAddress) {
        localStorage.setItem("userWalletAddress", walletAddress);
      }
    }

    App.account = walletAddress;
  },

  handleLogout: function() {
    localStorage.removeItem("userWalletAddress");

    App.login();
    App.markAdopted();
  }
};

$(function() {
  $(window).load(function() {
    App.init();
  });
});
