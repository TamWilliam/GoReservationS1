# Projet réservation salon de coiffure

Le but est de créer une solution de réservation pour les salons de coiffure un peu à la façon de Planity sur la technologie Go
Il y 3 types d’utilisateurs :

- Les salons de coiffure qui peuvent donner leurs créneaux d’ouverture (par salon et par coiffeur)
- Les clients qui, après s’être connectés, peuvent réserver un créneau
- Les administrateurs

## Objectifs

### Contexte

Répertoire de salons de coiffure situés à Paris

### Système de session

Créer un système de session : utilisateur, salon, admin

### Formulaire, datepicker

Le salon de coiffure doit pouvoir mettre ses créneaux
Créer un outil où, les utilisateurs connectés peuvent choisir un salon puis un créneau qui a été défini par le salon

### Front

Utilisation de Bootstrap pour le style
**page d'accueil**
**page login**
**page de création de compte** création d'un compte utilisateur
**page de compte utilisateur** update/delete d'un compte utilisateur
**page d'ajout de salon de coiffure** création d'un compte salon de coiffure
**page de compte salon de coiffure** update/delete d'un compte salon de coiffure
**page des résultats des salons** read des salons de coiffure
**page des réservations**
**page admin** read/delete comptes utilisateurs et salons de coiffure

### DB

- **customers**

  - id_customer _INT_
  - role _INT_ (0 : user, 1 : admin)
  - first_name _VARCHAR_
  - last_name _VARCHAR_
  - email _VARCHAR_
  - password _VARCHAR_

- **hair_salons**

  - id_hair_salon _INT_
  - name _VARCHAR_
  - address _VARCHAR_
  - email _VARCHAR_
  - password _VARCHAR

- **openinghours**

  - id_opening_hours _INT_
  - id_hair_salon _INT_
  - day _INT_
  - opening _TIME_
  - closing _TIME_

- **hair_dressers**

  - id_hair_dresser _INT_
  - id_hair_salon _INT_
  - first_name _VARCHAR_
  - last_name _VARCHAR_

- **hair_dresser_schedules**

  - id_hair_dresser_schedule _INT_
  - id_hair_dresser _INT_
  - day _INT_
  - start_shift _TIME_
  - end_shift _TIME_

- **reservations**
  - id_reservation _INT_
  - id_customer _INT_
  - id_hair_salon _INT_
  - id_hair_dresser _INT_
  - reservation_date _TIMESTAMP_
