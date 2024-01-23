# Інженерні думки: Promotion Aplication

*Стенографія відео ["#100 Інженерні думки | Vault Operator, Application CRD, Flux-TF controller"](https://www.youtu.be/zf9FScRhtCY) від Дениса Васильєва*

[YouTube TS](https://www.youtu.be/zf9FScRhtCY&t=990s)

Зрозуміло, що в нас є Helm Releases, ми можемо робити commit в repository далі Flux буде підтягувати це все і здійснювати реконсиляцію, якщо там будуть якісь зміни.

Але завданням є створити інструмент promotion application.

[YouTube TS](https://youtu.be/zf9FScRhtCY?t=794)

**Розлянемо реальний приклад.** В нас є Slack і на його базі за допомогою команд ми можемо управляти нашим promotion, як в перед, так і назад, тобто здійснювати promote і rollback.

[YouTube TS](https://youtu.be/zf9FScRhtCY?t=810)

**Пропонується реалізувати наступний підхід.**

GoLang application який нативно за допомогою бібліотек і модулів робить API calls в Kubernetes, використовує вбудовану SQLite базу данних і використовує API Aplication CRD (наприклад [Wordpress Application](https://github.com/kubernetes-sigs/application/blob/master/docs/examples/wordpress/application.yaml)) який розширює Kubernetes на об'єкт Application.

[YouTube TS](https://youtu.be/zf9FScRhtCY?t=852)

В чому як би суть Application CRD? В тому що в нас є Pods, StatefullSets, Deployments і таке інше. Але в нас з'явився ще один об'єкт, який може містити metadata, які безпосердньо описують наш bundle або додаток і оновлюється (updates) по статусу додатка в real time, тобто якщо щось змінилося в додатку ми бачимо цей статус в єдиному CRD об'єкті.

[YouTube TS](https://youtu.be/zf9FScRhtCY?t=891)

Суть в принциві проста. В нас є GoLang додаток, що реалізує бота на Slack і має connect до API Kubernetes, ми зчитуємо всі Application CRDs, які є на environments `dev`, `qa`, `stage`, `prod` (*це можуть бути різні namespaces, різні clusters і т.д.*) у вигляді списків (*дана операція відбувається дуже швидко, до секунди часу парсинг яких вже готовий*) які в свою чергу розкладуються в якусь inmemory embedded RDBMS. Після чого ми можемо робити різні operations, любий аналіз, любу калькуляцію і основне це promotion. Тобто, ми бачимо версії які змінено, і за допомогою певних команд зробити виборку версії додатку і зробити promotion певної версії на потрібний environment. Що при цьому робить наш Slack bot? Він використовуючи бібліотеку/модуль для роботи з OCI Registry змінює tag або переміщає сам image і далі в flow процесу підключається Flux. Flux відстежує OCI Registry і якщо змінився tag або версія image він оновлює helm release, для цього нам не потрібно заходити в репозиторій і змінювати щось.

Таким чином ми реалізуємол GitOps.

---

OCI, або Open Container Initiative, є ініціативою, спрямованою на розробку відкритого стандарту для контейнерів та інструментів для їх управління. 

OCI Registry - це сховище, яке дотримується стандартів OCI та використовує OCI-сумісні зображення контейнерів. 
