# Projet réservation salon de coiffure
Le but est de créer une solution de réservation pour les salons de coiffure un peu à la façon de Planity sur la technologie Go
Il y 3 types d’utilisateurs :
  – Les salons de coiffure qui peuvent donner leurs créneaux d’ouverture (par
    salon et par coiffeur)
  – Les clients qui, après s’être connectés, peuvent réserver un créneau
  – Les administrateurs

# Install mux
go get -u github.com/gorilla/mux

# Install pq
go get -u github.com/lib/pq