PGDMP     #    7                |            projetgoreservation    15.2    15.2 "    %           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                      false            &           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                      false            '           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                      false            (           1262    57484    projetgoreservation    DATABASE     �   CREATE DATABASE projetgoreservation WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'French_France.1252';
 #   DROP DATABASE projetgoreservation;
                postgres    false            �            1259    57485 	   customers    TABLE     �   CREATE TABLE public.customers (
    id_customer integer NOT NULL,
    role integer,
    first_name character varying(100),
    last_name character varying(100),
    email text,
    password character varying(100)
);
    DROP TABLE public.customers;
       public         heap    postgres    false            �            1259    57536    hair_dresser_schedules    TABLE     �   CREATE TABLE public.hair_dresser_schedules (
    id_hair_dresser_schedule integer NOT NULL,
    id_hair_dresser integer,
    day integer,
    start_shift time without time zone,
    end_shift time without time zone
);
 *   DROP TABLE public.hair_dresser_schedules;
       public         heap    postgres    false            �            1259    57510    hair_dressers    TABLE     �   CREATE TABLE public.hair_dressers (
    id_hair_dresser integer NOT NULL,
    first_name text,
    last_name character varying,
    id_hair_salon integer
);
 !   DROP TABLE public.hair_dressers;
       public         heap    postgres    false            �            1259    57490    hair_salons    TABLE     �   CREATE TABLE public.hair_salons (
    id_hair_salon integer NOT NULL,
    name text,
    address text,
    email character varying,
    password character varying(100)
);
    DROP TABLE public.hair_salons;
       public         heap    postgres    false            �            1259    57531    opening_hours    TABLE     �   CREATE TABLE public.opening_hours (
    id_opening_hours integer NOT NULL,
    id_hair_salon integer,
    day integer,
    opening time without time zone,
    closing time without time zone
);
 !   DROP TABLE public.opening_hours;
       public         heap    postgres    false            �            1259    57557    reservations    TABLE     �   CREATE TABLE public.reservations (
    id_reservation integer NOT NULL,
    id_customer integer,
    id_hair_salon integer,
    id_hair_dresser integer,
    reservation_date timestamp without time zone
);
     DROP TABLE public.reservations;
       public         heap    postgres    false            �            1259    57556    reservation_id_reservation_seq    SEQUENCE     �   CREATE SEQUENCE public.reservation_id_reservation_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 5   DROP SEQUENCE public.reservation_id_reservation_seq;
       public          postgres    false    220            )           0    0    reservation_id_reservation_seq    SEQUENCE OWNED BY     b   ALTER SEQUENCE public.reservation_id_reservation_seq OWNED BY public.reservations.id_reservation;
          public          postgres    false    219            y           2604    57560    reservations id_reservation    DEFAULT     �   ALTER TABLE ONLY public.reservations ALTER COLUMN id_reservation SET DEFAULT nextval('public.reservation_id_reservation_seq'::regclass);
 J   ALTER TABLE public.reservations ALTER COLUMN id_reservation DROP DEFAULT;
       public          postgres    false    220    219    220                      0    57485 	   customers 
   TABLE DATA           ^   COPY public.customers (id_customer, role, first_name, last_name, email, password) FROM stdin;
    public          postgres    false    214   �,                  0    57536    hair_dresser_schedules 
   TABLE DATA           x   COPY public.hair_dresser_schedules (id_hair_dresser_schedule, id_hair_dresser, day, start_shift, end_shift) FROM stdin;
    public          postgres    false    218   Z.                 0    57510    hair_dressers 
   TABLE DATA           ^   COPY public.hair_dressers (id_hair_dresser, first_name, last_name, id_hair_salon) FROM stdin;
    public          postgres    false    216   �.                 0    57490    hair_salons 
   TABLE DATA           T   COPY public.hair_salons (id_hair_salon, name, address, email, password) FROM stdin;
    public          postgres    false    215   �/                 0    57531    opening_hours 
   TABLE DATA           _   COPY public.opening_hours (id_opening_hours, id_hair_salon, day, opening, closing) FROM stdin;
    public          postgres    false    217   #2       "          0    57557    reservations 
   TABLE DATA           u   COPY public.reservations (id_reservation, id_customer, id_hair_salon, id_hair_dresser, reservation_date) FROM stdin;
    public          postgres    false    220   ~2       *           0    0    reservation_id_reservation_seq    SEQUENCE SET     M   SELECT pg_catalog.setval('public.reservation_id_reservation_seq', 11, true);
          public          postgres    false    219            }           2606    57496    hair_salons Hairdress_pkey 
   CONSTRAINT     e   ALTER TABLE ONLY public.hair_salons
    ADD CONSTRAINT "Hairdress_pkey" PRIMARY KEY (id_hair_salon);
 F   ALTER TABLE ONLY public.hair_salons DROP CONSTRAINT "Hairdress_pkey";
       public            postgres    false    215                       2606    57516    hair_dressers hairdresser_pkey 
   CONSTRAINT     i   ALTER TABLE ONLY public.hair_dressers
    ADD CONSTRAINT hairdresser_pkey PRIMARY KEY (id_hair_dresser);
 H   ALTER TABLE ONLY public.hair_dressers DROP CONSTRAINT hairdresser_pkey;
       public            postgres    false    216            �           2606    57540 /   hair_dresser_schedules hairdresserschedule_pkey 
   CONSTRAINT     �   ALTER TABLE ONLY public.hair_dresser_schedules
    ADD CONSTRAINT hairdresserschedule_pkey PRIMARY KEY (id_hair_dresser_schedule);
 Y   ALTER TABLE ONLY public.hair_dresser_schedules DROP CONSTRAINT hairdresserschedule_pkey;
       public            postgres    false    218            �           2606    57535    opening_hours openinghours_pkey 
   CONSTRAINT     k   ALTER TABLE ONLY public.opening_hours
    ADD CONSTRAINT openinghours_pkey PRIMARY KEY (id_opening_hours);
 I   ALTER TABLE ONLY public.opening_hours DROP CONSTRAINT openinghours_pkey;
       public            postgres    false    217            �           2606    57562    reservations reservation_pkey 
   CONSTRAINT     g   ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservation_pkey PRIMARY KEY (id_reservation);
 G   ALTER TABLE ONLY public.reservations DROP CONSTRAINT reservation_pkey;
       public            postgres    false    220            {           2606    57489    customers users_pkey 
   CONSTRAINT     [   ALTER TABLE ONLY public.customers
    ADD CONSTRAINT users_pkey PRIMARY KEY (id_customer);
 >   ALTER TABLE ONLY public.customers DROP CONSTRAINT users_pkey;
       public            postgres    false    214            �           2606    57541 &   hair_dressers fk_hairdresser_hairsalon    FK CONSTRAINT     �   ALTER TABLE ONLY public.hair_dressers
    ADD CONSTRAINT fk_hairdresser_hairsalon FOREIGN KEY (id_hair_salon) REFERENCES public.hair_salons(id_hair_salon);
 P   ALTER TABLE ONLY public.hair_dressers DROP CONSTRAINT fk_hairdresser_hairsalon;
       public          postgres    false    215    216    3197            �           2606    57546 9   hair_dresser_schedules fk_hairdresserschedule_hairdresser    FK CONSTRAINT     �   ALTER TABLE ONLY public.hair_dresser_schedules
    ADD CONSTRAINT fk_hairdresserschedule_hairdresser FOREIGN KEY (id_hair_dresser) REFERENCES public.hair_dressers(id_hair_dresser);
 c   ALTER TABLE ONLY public.hair_dresser_schedules DROP CONSTRAINT fk_hairdresserschedule_hairdresser;
       public          postgres    false    218    3199    216            �           2606    57551 '   opening_hours fk_openinghours_hairsalon    FK CONSTRAINT     �   ALTER TABLE ONLY public.opening_hours
    ADD CONSTRAINT fk_openinghours_hairsalon FOREIGN KEY (id_hair_salon) REFERENCES public.hair_salons(id_hair_salon);
 Q   ALTER TABLE ONLY public.opening_hours DROP CONSTRAINT fk_openinghours_hairsalon;
       public          postgres    false    217    215    3197            �           2606    57583 '   reservations fk_reservation_hairdresser    FK CONSTRAINT     �   ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT fk_reservation_hairdresser FOREIGN KEY (id_hair_dresser) REFERENCES public.hair_dressers(id_hair_dresser);
 Q   ALTER TABLE ONLY public.reservations DROP CONSTRAINT fk_reservation_hairdresser;
       public          postgres    false    220    3199    216            �           2606    57578 %   reservations fk_reservation_hairsalon    FK CONSTRAINT     �   ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT fk_reservation_hairsalon FOREIGN KEY (id_hair_salon) REFERENCES public.hair_salons(id_hair_salon);
 O   ALTER TABLE ONLY public.reservations DROP CONSTRAINT fk_reservation_hairsalon;
       public          postgres    false    215    220    3197            �           2606    57563 )   reservations reservation_id_customer_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservation_id_customer_fkey FOREIGN KEY (id_customer) REFERENCES public.customers(id_customer);
 S   ALTER TABLE ONLY public.reservations DROP CONSTRAINT reservation_id_customer_fkey;
       public          postgres    false    220    214    3195            �           2606    57573 ,   reservations reservation_id_hairdresser_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservation_id_hairdresser_fkey FOREIGN KEY (id_hair_dresser) REFERENCES public.hair_dressers(id_hair_dresser);
 V   ALTER TABLE ONLY public.reservations DROP CONSTRAINT reservation_id_hairdresser_fkey;
       public          postgres    false    3199    220    216            �           2606    57568 *   reservations reservation_id_hairsalon_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservation_id_hairsalon_fkey FOREIGN KEY (id_hair_salon) REFERENCES public.hair_salons(id_hair_salon);
 T   ALTER TABLE ONLY public.reservations DROP CONSTRAINT reservation_id_hairsalon_fkey;
       public          postgres    false    3197    220    215               �  x�U����0Ek�c_"��#�F��&�p8���@6l��3�n� ,永8sIPF���Y�vV'���/8��E�d)T5_(V���0'��]�V�z�3��e��I�᲍��䨔��`���8�cr�k'O�D�p�����'���ǖ�7&S1.�'��w�-#e����k��	��*��/�o1!SJ���Ħ�ͧJ8���8}��z^�����C٦]D+�3���Fc�N�\�����Q\^�����(�����0���$��B�=y�FҖ+�����<�&TO����R��_��4H���1�f�w��o��$n�x����羪��p}��pLp5�VS��\��;�Q"X�:��{�s�C�+���e��Qn�e���f�ȍ�ۛB�`D�����e��C-�4�F�,��]�h�ML({�@���)�߃���>ݝ          n   x�m�� !D�P���ֲ�ױc���F�9y�J�#k����\@������+h%��G�x;yɼ{:ҝT�G�qzkz�,v�-����%�`��S�N��1xf� ��(�           x���n�@E�w>�b�{ɫPZZZ��t�&#�b2�� j��N���J���,qp�����\������},"i���`ɔ	s�{6bU�I+�'���Fxd��"�^��X�s�I�[���B�1JM��)�(�ĩ[;t��mh.	;��x�6���Qn0r��S����[Q¡�����^�"�t���o�o�ؑ��1q~�=����\�
����9cMZZ|��TƐ�!��q����@���B����{���h�	��2�������gm           x�e�Mn�0���)x����w��h�EI�U7Crh��B?.|���+e����H���7O�7솞Ǳ��T�_�x���p���a�V-���iuߥ�ȰԨ��P��E�Y8	�ȆF��];�v�o�ہic��}��D3��q��0e��/e�.��
�T���W0�ݡM|��s>�o8������T�T�d��R�1��BHZ2I�h���qX�A��$��O��pt���E�!kr��M�(��<�%���(�=5��@{��K�L	ŷ��?cO5�M�@o��u*e��"I����.�b�ȩ��e�4���H}�؎���ǋ>_�[��Yj"Y)���5z��4��uD�ꞇ���B���u�̦������Rv��"Sh*�(��)N��y��p����cbF���3�?KG'3K��[-7t�H�`	�֩�����X�*�&�W��륹gf����{�y��OENɓ�E{Y�W@iR����2��=���'��0�����>��E� �M��)����{㛟wM��y4�         K   x�m���0�7SHB�Y���D�~�yX��&��> OlWQ��`�(�*�h|�G�>w|�}5�I�ƟKU_�>�      "      x������ � �     