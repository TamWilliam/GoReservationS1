# Projet réservation salon de coiffure
Le but est de créer une solution de réservation pour les salons de coiffure un peu à la façon de Planity sur la technologie Go
Il y 3 types d’utilisateurs :
  – Les salons de coiffure qui peuvent donner leurs créneaux d’ouverture (par salon et par coiffeur)
  – Les clients qui, après s’être connectés, peuvent réserver un créneau
  – Les administrateurs

## Objectifs
### Contexte
Répertoire de salons de coiffure situés à Paris

### Système de session
Créer un système de session : utilisateur, salon, admin

### Formulaire, datepicker
Le salon de coiffure doit pouvoir mettre ses créneaux
Créer un outil où, les utilisateurs connectés peuvent choisir un salon puis un créneau qui a été défini par le salon

### Front
Utilisation de Tailwind CSS / Bootstrap
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
- **Users**
  - id_user *INT*
  - role *INT* (0 : user, 1 : admin)
  - firstname *VARCHAR*
  - lastname *VARCHAR*
  - email *VARCHAR*
  - password *VARCHAR*

- **Hairdress**
  - id_hairdress *INT*
  - id_haidresser *INT*
  - employees *INT* (le nombre d'employés)
  - name *VARCHAR*
  - address *TEXT*
  - email *VARCHAR*
  - hours *VARCHAR*
  - days *VARCHAR*
  - password *VARCHAR*

- **Hairdresser**
  - id_hairdresser *INT*
  - id_hairdress *INT*
  - firstname *VARCHAR*
  - lastname *VARCHAR*
  - hours *VARCHAR*
  - days *VARCHAR*