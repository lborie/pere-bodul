# Père Bodul

Père Bodul est une application web simple qui génère des histoires pour enfants en utilisant l'API ChatGPT.

## Fonctionnalités

- Permet de spécifier plusieurs paramètres pour personnaliser l'histoire :
    - L'âge de l'enfant
    - Un héros
    - Un méchant
    - Le lieu de l'histoire
    - Des objets à incorporer dans l'histoire

## Prérequis

Pour faire fonctionner l'application, vous avez besoin de :

- Go (version 1.16 ou supérieure)
- Une clé API pour l'API OpenAI

## Installation

1. Clonez ce dépôt sur votre machine locale.
2. Installez les dépendances en exécutant `go get` dans le répertoire du projet.
3. Configurez la clé API OpenAI en définissant la variable d'environnement `OPENAI_KEY`.

## Utilisation

1. Exécutez le serveur en lançant `go run main.go`.
2. Ouvrez un navigateur web et accédez à `http://localhost:8080`.
3. Remplissez le formulaire avec vos préférences pour l'histoire et cliquez sur "Générer l'histoire".
4. L'histoire générée s'affichera sous le formulaire.

## Licence

Ce projet est sous licence MIT. Pour plus de détails, voir le fichier `LICENSE`.
